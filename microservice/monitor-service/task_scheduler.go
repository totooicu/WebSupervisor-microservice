package main

import (
	"log"
	"sync"
	"time"
)

// TaskScheduler 定时任务调度器
type TaskScheduler struct {
	jobs           []JobConfig
	communicator   *ServiceCommunicator
	parser         *ConfigParser
	debug          bool
	tickerMap      map[string]*time.Ticker
	tickerMutex    sync.Mutex
	stopChan       chan struct{}
}

// NewTaskScheduler 创建定时任务调度器
func NewTaskScheduler(jobs []JobConfig, communicator *ServiceCommunicator, parser *ConfigParser, debug bool) *TaskScheduler {
	return &TaskScheduler{
		jobs:           jobs,
		communicator:   communicator,
		parser:         parser,
		debug:          debug,
		tickerMap:      make(map[string]*time.Ticker),
		stopChan:       make(chan struct{}),
	}
}

// Start 启动定时任务调度器
func (ts *TaskScheduler) Start() {
	log.Println("Starting task scheduler...")
	
	// 为每个job启动定时任务
	for _, job := range ts.jobs {
		ts.startJobScheduler(job)
	}
	
	log.Println("Task scheduler started successfully")
}

// startJobScheduler 启动单个job的定时调度
func (ts *TaskScheduler) startJobScheduler(job JobConfig) {
	jobKey := job.ProjectName
	
	ts.tickerMutex.Lock()
	defer ts.tickerMutex.Unlock()
	
	// 检查是否已经有相同的job在运行
	if _, exists := ts.tickerMap[jobKey]; exists {
		log.Printf("Job %s is already running, skipping...", jobKey)
		return
	}
	
	// 创建定时器
	ticker := time.NewTicker(time.Duration(job.IntervalSecond) * time.Second)
	ts.tickerMap[jobKey] = ticker
	
	go func(job JobConfig, ticker *time.Ticker) {
		log.Printf("Starting scheduler for job: %s, interval: %d seconds", job.ProjectName, job.IntervalSecond)
		
		// 立即执行一次
		ts.executeJob(job)
		
		for {
			select {
			case <-ticker.C:
				ts.executeJob(job)
			case <-ts.stopChan:
				log.Printf("Stopping scheduler for job: %s", job.ProjectName)
				ticker.Stop()
				return
			}
		}
	}(job, ticker)
}

// executeJob 执行单个job
func (ts *TaskScheduler) executeJob(job JobConfig) {
	log.Printf("Executing job: %s", job.ProjectName)
	
	var wg sync.WaitGroup
	for _, urlConfig := range job.URLs {
		wg.Add(1)
		go func(urlConfig URLConfig) {
			defer wg.Done()
			if err := ts.executeURLTask(urlConfig); err != nil {
				log.Printf("Error executing URL task %s: %v", urlConfig.URL, err)
				
				// 发送错误通知
				errorMsg := "Failed to execute URL task: " + urlConfig.URL
				if notifyErr := ts.communicator.SendNotification("error", errorMsg); notifyErr != nil {
					log.Printf("Failed to send error notification: %v", notifyErr)
				}
			}
		}(urlConfig)
	}
	
	wg.Wait()
	log.Printf("Job %s execution completed", job.ProjectName)
}

// executeURLTask 执行单个URL任务
func (ts *TaskScheduler) executeURLTask(urlConfig URLConfig) error {
	if ts.debug {
		log.Printf("Debug - Executing URL task: %s", urlConfig.URL)
	}
	
	// 发送任务到爬虫服务
	if err := ts.communicator.SendToCrawlerService(urlConfig); err != nil {
		return err
	}
	
	return nil
}

// Stop 停止定时任务调度器
func (ts *TaskScheduler) Stop() {
	log.Println("Stopping task scheduler...")
	
	// 关闭stop通道
	close(ts.stopChan)
	
	// 停止所有定时器
	ts.tickerMutex.Lock()
	defer ts.tickerMutex.Unlock()
	
	for jobKey, ticker := range ts.tickerMap {
		ticker.Stop()
		log.Printf("Stopped ticker for job: %s", jobKey)
	}
	
	// 清空定时器映射
	ts.tickerMap = make(map[string]*time.Ticker)
	
	log.Println("Task scheduler stopped successfully")
}

// AddJob 添加新的job到调度器
func (ts *TaskScheduler) AddJob(job JobConfig) error {
	if ts.debug {
		log.Printf("Debug - Adding job: %s", job.ProjectName)
	}
	
	// 验证job配置
	if err := ts.parser.validateJobsConfig(&job); err != nil {
		return err
	}
	
	// 添加到jobs列表
	ts.jobs = append(ts.jobs, job)
	
	// 启动该job的调度
	ts.startJobScheduler(job)
	
	log.Printf("Job %s added successfully", job.ProjectName)
	return nil
}

// RemoveJob 从调度器移除job
func (ts *TaskScheduler) RemoveJob(projectName string) {
	if ts.debug {
		log.Printf("Debug - Removing job: %s", projectName)
	}
	
	ts.tickerMutex.Lock()
	defer ts.tickerMutex.Unlock()
	
	// 停止定时器
	if ticker, exists := ts.tickerMap[projectName]; exists {
		ticker.Stop()
		delete(ts.tickerMap, projectName)
		log.Printf("Removed ticker for job: %s", projectName)
	}
	
	// 从jobs列表中移除
	newJobs := []JobConfig{}
	for _, job := range ts.jobs {
		if job.ProjectName != projectName {
			newJobs = append(newJobs, job)
		}
	}
	ts.jobs = newJobs
	
	log.Printf("Job %s removed successfully", projectName)
}

// UpdateJob 更新现有job
func (ts *TaskScheduler) UpdateJob(job JobConfig) error {
	if ts.debug {
		log.Printf("Debug - Updating job: %s", job.ProjectName)
	}
	
	// 验证job配置
	if err := ts.parser.validateJobsConfig(&job); err != nil {
		return err
	}
	
	// 先移除旧的job
	ts.RemoveJob(job.ProjectName)
	
	// 添加新的job
	return ts.AddJob(job)
}