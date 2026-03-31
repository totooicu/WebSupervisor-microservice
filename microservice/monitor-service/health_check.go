package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"WebSupervisor/MyTool/streamtool"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	config      *MonitorConfig
	debug       bool
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(config *MonitorConfig, debug bool) *HealthChecker {
	return &HealthChecker{
		config:      config,
		debug:       debug,
	}
}

// StartHealthCheckServer 启动健康检查服务器
func (hc *HealthChecker) StartHealthCheckServer() {
	http.HandleFunc("/health", hc.healthCheckHandler)
	
	server := &http.Server{
		Addr:         ":" + hc.config.HealthCheckPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	log.Printf("Starting health check server on port %s", hc.config.HealthCheckPort)
	
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start health check server: %v", err)
		}
	}()
}

// healthCheckHandler 健康检查处理函数
func (hc *HealthChecker) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := hc.checkHealth()
	
	if status.Healthy {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "unhealthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `", "issues": "` + status.Issues + `"}`))
	}
}

// HealthStatus 健康状态
type HealthStatus struct {
	Healthy bool
	Issues  string
}

// checkHealth 检查服务健康状态
func (hc *HealthChecker) checkHealth() HealthStatus {
	status := HealthStatus{
		Healthy: true,
		Issues:  "",
	}
	
	// 检查Redis连接（通过streamtool）
	if err := hc.checkRedisConnection(); err != nil {
		status.Healthy = false
		status.Issues += "Redis connection failed: " + err.Error() + "; "
	}
	
	// 检查配置文件
	if err := hc.checkConfigFiles(); err != nil {
		status.Healthy = false
		status.Issues += "Config files check failed: " + err.Error() + "; "
	}
	
	return status
}

// checkRedisConnection 检查Redis连接
func (hc *HealthChecker) checkRedisConnection() error {
	// 创建一个简单的测试消息来验证Redis连接
	testMsg := map[string]interface{}{
		"test": "health_check",
		"time": time.Now().Format(time.RFC3339),
	}
	
	// 尝试发布消息到测试流
	testStream := "health-check-test"
	st := streamtool.GetStreamTool()
	if !st.StreamPush(testMsg, testStream) {
		return fmt.Errorf("Redis connection test failed")
	}
	
	return nil
}

// checkConfigFiles 检查配置文件
func (hc *HealthChecker) checkConfigFiles() error {
	// 检查jobs配置文件是否存在
	if _, err := os.Stat(hc.config.JobsPath); os.IsNotExist(err) {
		return fmt.Errorf("jobs config file not found: %s", hc.config.JobsPath)
	}
	
	return nil
}