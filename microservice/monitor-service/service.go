package main

import (
	"log"
	"time"

	"WebSupervisor/model"
	"WebSupervisor/MyTool/redis"
)

// JobConfig 任务配置结构
type JobConfig struct {
	ProjectName    string            `json:"projectName"`
	Header         map[string]string `json:"header"`
	URLs           []URLConfig       `json:"urls"`
	IntervalSecond int               `json:"intervalSecond"`
}

// URLConfig URL配置结构
type URLConfig struct {
	URL            string            `json:"url"`
	Output         string            `json:"output"`
	Test           bool              `json:"test"`
	Header         map[string]string `json:"header"`
	Method         string            `json:"method"`
	Body           map[string]string `json:"body"`
	StringPlayLoad string            `json:"stringPlayLoad"`
	Type           string            `json:"type"`
	JSONKeys       []JSONKey         `json:"jsonKeys"`
	HTMLKeys       []HTMLKey         `json:"htmlKeys"`
	IntervalSecond int               `json:"intervalSecond"`
}

// JSONKey JSON解析配置
type JSONKey struct {
	Path []interface{} `json:"path"`
	Keys []string      `json:"key"`
}

// HTMLKey HTML解析配置
type HTMLKey struct {
	Left  string   `json:"left"`
	Right string   `json:"right"`
	Keys  []string `json:"key"`
}

// MonitorService 监控服务
type MonitorService struct {
	redisClient   *redis.Client
	config        *model.MonitorConfig
	parser        *ConfigParser
	communicator  *ServiceCommunicator
	scheduler     *TaskScheduler
	healthChecker *HealthChecker
	debug         bool
}

// NewMonitorService 创建监控服务实例
func NewMonitorService(config *model.MonitorConfig, debug bool) *MonitorService {
	redisClient := redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB)
	
	// 创建配置解析器
	parser := NewConfigParser(debug)
	
	// 创建服务通信器
	communicator := NewServiceCommunicator(redisClient, config, debug)
	
	// 创建健康检查器
	healthChecker := NewHealthChecker(redisClient, config, debug)
	
	return &MonitorService{
		redisClient:   redisClient,
		config:        config,
		parser:        parser,
		communicator:  communicator,
		healthChecker: healthChecker,
		debug:         debug,
	}
}

// Start 启动监控服务
func (s *MonitorService) Start() {
	log.Println("Starting monitor service...")
	
	if s.debug {
		log.Printf("Debug mode enabled, config: %+v", s.config)
	}
	
	// 启动健康检查服务器
	s.healthChecker.StartHealthCheckServer()
	
	// 加载并解析jobs配置
	jobs, err := s.loadJobs()
	if err != nil {
		log.Fatalf("Failed to load jobs: %v", err)
	}
	
	// 创建任务调度器
	s.scheduler = NewTaskScheduler(jobs, s.communicator, s.parser, s.debug)
	
	// 启动事件监听器
	s.communicator.StartEventListener()
	
	// 启动任务调度器
	s.scheduler.Start()
	
	log.Println("Monitor service started successfully")
	
	// 保持服务运行
	select {}
}

// loadJobs 加载并解析jobs配置
func (s *MonitorService) loadJobs() ([]JobConfig, error) {
	if s.debug {
		log.Printf("Debug - Loading jobs from: %s", s.config.JobsPath)
	}
	
	// 解析jobs配置文件
	jobConfig, err := s.parser.ParseJobsConfig(s.config.JobsPath)
	if err != nil {
		return nil, err
	}
	
	log.Printf("Loaded job: %s with %d URLs", jobConfig.ProjectName, len(jobConfig.URLs))
	
	return []JobConfig{*jobConfig}, nil
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return "task-" + time.Now().Format("20060102150405")
}