package main

import (
	"fmt"
	"log"

	"WebSupervisor/MyTool"
)

// ConfigParser 配置解析器
type ConfigParser struct {
	debug bool
}

// NewConfigParser 创建配置解析器
func NewConfigParser(debug bool) *ConfigParser {
	return &ConfigParser{
		debug: debug,
	}
}

// ParseJobsConfig 解析jobs配置文件
func (p *ConfigParser) ParseJobsConfig(configPath string) (*JobConfig, error) {
	if p.debug {
		log.Printf("Debug - Parsing jobs config from: %s", configPath)
	}

	var jobConfig JobConfig
	if err := MyTool.LoadConfig(configPath, &jobConfig); err != nil {
		return nil, fmt.Errorf("failed to load jobs config: %w", err)
	}

	// 验证配置参数
	if err := p.validateJobsConfig(&jobConfig); err != nil {
		return nil, fmt.Errorf("invalid jobs config: %w", err)
	}

	if p.debug {
		log.Printf("Debug - Jobs config parsed successfully: %+v", jobConfig)
	}

	return &jobConfig, nil
}

// validateJobsConfig 验证jobs配置参数
func (p *ConfigParser) validateJobsConfig(config *JobConfig) error {
	// 验证项目名称
	if config.ProjectName == "" {
		return fmt.Errorf("projectName is required")
	}

	// 验证间隔时间
	if config.IntervalSecond <= 0 {
		return fmt.Errorf("intervalSecond must be positive")
	}

	// 验证URLs列表
	if len(config.URLs) == 0 {
		return fmt.Errorf("urls list cannot be empty")
	}

	// 验证每个URL配置
	for i, urlConfig := range config.URLs {
		if err := p.validateURLConfig(&urlConfig); err != nil {
			return fmt.Errorf("url %d validation failed: %w", i, err)
		}
	}

	return nil
}

// validateURLConfig 验证URL配置参数
func (p *ConfigParser) validateURLConfig(config *URLConfig) error {
	// 验证URL
	if config.URL == "" {
		return fmt.Errorf("url is required")
	}

	// 验证方法（默认GET）
	if config.Method == "" {
		config.Method = "GET"
	}

	// 验证类型（默认html）
	if config.Type == "" {
		config.Type = "html"
	}

	// 验证间隔时间（使用job级别的默认值）
	if config.IntervalSecond <= 0 {
		config.IntervalSecond = 60 // 默认60秒
	}

	return nil
}

// ParseEmailConfig 解析邮箱配置（从MonitorConfig中提取）
func (p *ConfigParser) ParseEmailConfig(emailConfig map[string]string) (map[string]string, error) {
	if p.debug {
		log.Printf("Debug - Parsing email config: %v", emailConfig)
	}

	// 验证邮箱配置
	if emailConfig == nil {
		return nil, fmt.Errorf("email config is required")
	}

	requiredFields := []string{"smtp_server", "smtp_port", "username", "password", "from_email", "to_email"}
	for _, field := range requiredFields {
		if emailConfig[field] == "" {
			return nil, fmt.Errorf("email config missing required field: %s", field)
		}
	}

	if p.debug {
		log.Println("Debug - Email config validated successfully")
	}

	return emailConfig, nil
}