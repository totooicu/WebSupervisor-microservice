package main

import (
	"WebSupervisor/model"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

// TaskScheduler 定时任务调度器
type TaskScheduler struct {
	jobs         []JobConfig
	communicator *ServiceCommunicator
	parser       *ConfigParser
	debug        bool
	tickerMap    map[string]*time.Ticker
	tickerMutex  sync.Mutex
	stopChan     chan struct{}
}

// NewTaskScheduler 创建定时任务调度器
func NewTaskScheduler(jobs []JobConfig, communicator *ServiceCommunicator, parser *ConfigParser, debug bool) *TaskScheduler {
	return &TaskScheduler{
		jobs:         jobs,
		communicator: communicator,
		parser:       parser,
		debug:        debug,
		tickerMap:    make(map[string]*time.Ticker),
		stopChan:     make(chan struct{}),
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
		ts.executeJob(job) //即执行crawler-service

		for {
			select {
			case <-ticker.C: //定时器触发
				ts.executeJob(job)
			case <-ts.stopChan:
				log.Printf("Stopping scheduler for job: %s", job.ProjectName)
				ticker.Stop()
				return
			}
			//检查crawler-service_input_stream是否有新任务

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

// executeURLTask 执行单个URL任务，整合爬取、解析、比对缓存、通知功能
func (ts *TaskScheduler) executeURLTask(urlConfig URLConfig) error {
	if ts.debug {
		log.Printf("Debug - Executing URL task: %s", urlConfig.URL)
	}
	log.Printf(">>>Executing URL task: %V", urlConfig.URL)

	// 1. 发送爬虫请求，获取响应对象（非阻塞）
	crawlerResponseObj, err := ts.communicator.SendMessageWithResponse(
		ts.communicator.crawlerStream,
		ts.createCrawlerTask(urlConfig),
	)
	if err != nil {
		return fmt.Errorf("crawler failed: %w", err)
	}

	if ts.debug {
		log.Printf("Debug - Crawler request sent, waiting for response")
	}

	// 在这里可以执行一些与爬虫响应无关的计算任务
	// ...

	// 获取爬虫响应（阻塞）
	crawlerResponseData := crawlerResponseObj.Get()

	// 解析爬虫响应数据
	// var crawlerResult struct {
	// 	Content string `json:"content"`
	// }
	// if err := json.Unmarshal([]byte(crawlerResponseData.(string)), &crawlerResult); err != nil {
	// 	return fmt.Errorf("failed to parse crawler response: %w", err)
	// }

	crawlerResult := crawlerResponseData.(map[string]interface{})["content"].(string)
	// 2. 发送解析请求，获取响应对象（非阻塞）
	parserResponseObj, err := ts.communicator.SendMessageWithResponse(
		ts.communicator.parserStream,
		ts.createParserTask(crawlerResult, urlConfig),
	)
	if err != nil {
		return fmt.Errorf("parser failed: %w", err)
	}

	if ts.debug {
		log.Printf("Debug - Parser request sent, waiting for response")
	}

	// 在这里可以执行一些与解析响应无关的计算任务
	// ...

	// 获取解析响应（阻塞）
	parserResponseData := parserResponseObj.Get()

	// 解析解析器响应数据
	var parserResult struct {
		ParsedData interface{} `json:"parsed_data"`
	}
	if err := json.Unmarshal([]byte(parserResponseData.(string)), &parserResult); err != nil {
		return fmt.Errorf("failed to parse parser response: %w", err)
	}

	// 3. 发送缓存比对请求，获取响应对象（非阻塞）
	cacheResponseObj, err := ts.communicator.SendMessageWithResponse(
		ts.communicator.cacheStream,
		ts.createCacheTask(parserResult.ParsedData, urlConfig),
	)
	if err != nil {
		return fmt.Errorf("cache compare failed: %w", err)
	}

	if ts.debug {
		log.Printf("Debug - Cache compare request sent, waiting for response")
	}

	// 在这里可以执行一些与缓存比对响应无关的计算任务
	// ...

	// 获取缓存比对响应（阻塞）
	cacheResponseData := cacheResponseObj.Get()

	// 解析缓存响应数据
	var cacheResult struct {
		Changed bool `json:"changed"`
	}
	if err := json.Unmarshal([]byte(cacheResponseData.(string)), &cacheResult); err != nil {
		return fmt.Errorf("failed to parse cache response: %w", err)
	}

	// 4. 发送通知（如果数据有变化）
	if cacheResult.Changed {
		notificationMsg := fmt.Sprintf("Data changed for URL: %s", urlConfig.URL)
		if err := ts.communicator.SendNotification("info", notificationMsg); err != nil {
			log.Printf("Failed to send notification: %v", err)
		} else {
			log.Printf("Notification sent: %s", notificationMsg)
		}
	} else if ts.debug {
		log.Printf("Debug - No data change detected for URL: %s", urlConfig.URL)
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

func (ts *TaskScheduler) Fllow(urlConfig URLConfig) {
	ts.executeURLTask(urlConfig)
}

// createCrawlerTask 创建爬虫任务
func (ts *TaskScheduler) createCrawlerTask(urlConfig URLConfig) model.Message {
	taskID := strconv.Itoa(int(time.Now().UnixNano()))

	paramData := map[string]interface{}{
		"url":         urlConfig.URL,
		"method":      urlConfig.Method,
		"headers":     urlConfig.Header,
		"body":        urlConfig.Body,
		"str_payload": urlConfig.StringPlayLoad,
		"output":      urlConfig.Output,
		"test":        urlConfig.Test,
	}

	paramBytes, _ := json.Marshal(paramData)

	return model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "crawler-group",
		CallbackStream: "monitor-response-stream",
		ServiceName:    "http_request",
		Playload:       string(paramBytes),
	}
}

// createParserTask 创建解析任务
func (ts *TaskScheduler) createParserTask(rawData string, urlConfig URLConfig) model.Message {
	taskID := strconv.Itoa(int(time.Now().UnixNano()))

	paramData := map[string]interface{}{
		"content":   rawData,
		"type":      urlConfig.Type,
		"json_keys": urlConfig.JSONKeys,
		"html_keys": urlConfig.HTMLKeys,
		"output":    urlConfig.Output,
	}

	paramBytes, _ := json.Marshal(paramData)

	return model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "parser-group",
		CallbackStream: "monitor-response-stream",
		ServiceName:    "parse_data",
		Playload:       string(paramBytes),
	}
}

// createCacheTask 创建缓存比对任务
func (ts *TaskScheduler) createCacheTask(parsedData interface{}, urlConfig URLConfig) model.Message {
	taskID := strconv.Itoa(int(time.Now().UnixNano()))

	paramData := map[string]interface{}{
		"app":  "monitor-service",
		"key":  urlConfig.URL,
		"data": parsedData,
	}

	paramBytes, _ := json.Marshal(paramData)

	return model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "cache-group",
		CallbackStream: "monitor-response-stream",
		ServiceName:    "compare_and_save",
		Playload:       string(paramBytes),
	}
}
