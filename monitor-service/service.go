package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"WebSupervisor/new/model"
	"WebSupervisor/new/redis"
)

type MonitorService struct {
	redisClient *redis.Client
	config      *model.MonitorConfig
	jobs        []JobConfig
	debug       bool
}

type JobConfig struct {
	ProjectName    string            `json:"projectName"`
	Header         map[string]string `json:"header"`
	URLs           []URLConfig       `json:"urls"`
	IntervalSecond int               `json:"intervalSecond"`
}

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

type JSONKey struct {
	Path []interface{} `json:"path"`
	Keys []string      `json:"key"`
}

type HTMLKey struct {
	Left  string   `json:"left"`
	Right string   `json:"right"`
	Keys  []string `json:"key"`
}

func NewMonitorService(config *model.MonitorConfig, debug bool) *MonitorService {
	return &MonitorService{
		redisClient: redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB),
		config:      config,
		debug:       debug,
	}
}

func (s *MonitorService) Start() {
	log.Println("Starting monitor service...")
	if s.debug {
		log.Printf("Debug mode enabled, config: %+v", s.config)
	}

	if err := s.LoadJobs(s.config.JobsPath); err != nil {
		log.Fatalf("Failed to load jobs: %v", err)
	}

	for {
		if s.debug {
			log.Printf("Debug - Scheduling tasks, interval: %d seconds", s.config.IntervalSecond)
		}
		s.ScheduleTasks()
		time.Sleep(time.Second * time.Duration(s.config.IntervalSecond))
	}
}

func (s *MonitorService) LoadJobs(configPath string) error {
	if s.debug {
		log.Printf("Debug - Loading jobs from: %s", configPath)
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	var job JobConfig
	if err := json.Unmarshal(data, &job); err != nil {
		return err
	}

	s.jobs = append(s.jobs, job)
	log.Printf("Loaded %d jobs", len(s.jobs))

	if s.debug {
		log.Printf("Debug - Job details: %+v", job)
	}

	return nil
}

func (s *MonitorService) ScheduleTasks() {
	if s.debug {
		log.Printf("Debug - Scheduling %d jobs", len(s.jobs))
	}

	for _, job := range s.jobs {
		if s.debug {
			log.Printf("Debug - Processing job: %s, URLs: %d", job.ProjectName, len(job.URLs))
		}
		for _, urlConfig := range job.URLs {
			s.ScheduleTask(urlConfig)
		}
	}
}

func (s *MonitorService) ScheduleTask(urlConfig URLConfig) {
	taskID := generateTaskID()

	paramData := map[string]interface{}{
		"url":         urlConfig.URL,
		"method":      urlConfig.Method,
		"headers":     urlConfig.Header,
		"body":        urlConfig.Body,
		"str_payload": urlConfig.StringPlayLoad,
	}
	paramBytes, _ := json.Marshal(paramData)
	
	task := model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "parser-group",
		CallbackStream: "parser-tasks",
		ServiceName:    "http_request",
		Playload:       string(paramBytes),
	}

	log.Printf("Scheduling task: %s for URL: %s", taskID, urlConfig.URL)

	if s.debug {
		log.Printf("Debug - Task details: %+v", task)
	}

	if err := s.redisClient.PublishMessage(s.config.CrawlerServiceInputStream, task); err != nil {
		log.Printf("Error publishing task: %v", err)
	} else if s.debug {
		log.Printf("Debug - Published task to stream: %s", s.config.CrawlerServiceInputStream)
	}
}

func generateTaskID() string {
	return "task-" + time.Now().Format("20060102150405")
}