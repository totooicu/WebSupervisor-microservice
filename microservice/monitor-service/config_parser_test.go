package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfigParser_ParseJobsConfig 测试配置解析器解析jobs配置
func TestConfigParser_ParseJobsConfig(t *testing.T) {
	// 创建测试配置文件
	testConfig := `{
		"projectName": "test-project",
		"header": {},
		"intervalSecond": 60,
		"urls": [
			{
				"url": "https://example.com/api",
				"output": "./output.json",
				"test": false,
				"header": {},
				"method": "GET",
				"body": {},
				"type": "json",
				"jsonKeys": [],
				"htmlKeys": []
			}
		]
	}`

	// 写入临时文件
	tempFile, err := os.CreateTemp("", "jobs-test-*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	_, err = tempFile.WriteString(testConfig)
	assert.NoError(t, err)
	tempFile.Close()

	// 创建解析器并测试
	parser := NewConfigParser(false)
	jobConfig, err := parser.ParseJobsConfig(tempFile.Name())
	
	assert.NoError(t, err)
	assert.NotNil(t, jobConfig)
	assert.Equal(t, "test-project", jobConfig.ProjectName)
	assert.Equal(t, 60, jobConfig.IntervalSecond)
	assert.Len(t, jobConfig.URLs, 1)
	assert.Equal(t, "https://example.com/api", jobConfig.URLs[0].URL)
}

// TestConfigParser_ValidateJobsConfig 测试配置验证
func TestConfigParser_ValidateJobsConfig(t *testing.T) {
	parser := NewConfigParser(false)

	// 测试无效配置
	invalidConfig := &JobConfig{
		ProjectName:    "", // 空项目名
		IntervalSecond: 60,
		URLs: []URLConfig{
			{
				URL: "https://example.com",
			},
		},
	}

	err := parser.validateJobsConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "projectName is required")

	// 测试有效配置
	validConfig := &JobConfig{
		ProjectName:    "valid-project",
		IntervalSecond: 60,
		URLs: []URLConfig{
			{
				URL: "https://example.com",
			},
		},
	}

	err = parser.validateJobsConfig(validConfig)
	assert.NoError(t, err)
}

// TestConfigParser_ValidateURLConfig 测试URL配置验证
func TestConfigParser_ValidateURLConfig(t *testing.T) {
	parser := NewConfigParser(false)

	// 测试无效URL配置
	invalidURLConfig := &URLConfig{
		URL: "", // 空URL
	}

	err := parser.validateURLConfig(invalidURLConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "url is required")

	// 测试有效URL配置
	validURLConfig := &URLConfig{
		URL: "https://example.com",
	}

	err = parser.validateURLConfig(validURLConfig)
	assert.NoError(t, err)
	assert.Equal(t, "GET", validURLConfig.Method) // 验证默认值
	assert.Equal(t, "html", validURLConfig.Type) // 验证默认值
}

// TestConfigParser_ParseEmailConfig 测试邮箱配置解析
func TestConfigParser_ParseEmailConfig(t *testing.T) {
	parser := NewConfigParser(false)

	// 测试无效邮箱配置
	invalidEmailConfig := map[string]string{
		"smtp_server": "smtp.example.com",
		// 缺少必填字段
	}

	_, err := parser.ParseEmailConfig(invalidEmailConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	// 测试有效邮箱配置
	validEmailConfig := map[string]string{
		"smtp_server": "smtp.example.com",
		"smtp_port":   "587",
		"username":    "test@example.com",
		"password":    "password123",
		"from_email":  "test@example.com",
		"to_email":    "user@example.com",
	}

	result, err := parser.ParseEmailConfig(validEmailConfig)
	assert.NoError(t, err)
	assert.Equal(t, validEmailConfig, result)
}

// TestConfigParser_WithEnvironmentVariables 测试环境变量替换
func TestConfigParser_WithEnvironmentVariables(t *testing.T) {
	// 设置环境变量
	os.Setenv("TEST_URL", "https://test.example.com")
	defer os.Unsetenv("TEST_URL")

	// 创建包含环境变量的配置文件
	testConfig := `{
		"projectName": "test-project",
		"intervalSecond": 60,
		"urls": [
			{
				"url": "${TEST_URL}",
				"output": "./output.json",
				"method": "GET"
			}
		]
	}`

	// 写入临时文件
	tempFile, err := os.CreateTemp("", "jobs-env-*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	_, err = tempFile.WriteString(testConfig)
	assert.NoError(t, err)
	tempFile.Close()

	// 创建解析器并测试
	parser := NewConfigParser(false)
	jobConfig, err := parser.ParseJobsConfig(tempFile.Name())
	
	assert.NoError(t, err)
	assert.Equal(t, "https://test.example.com", jobConfig.URLs[0].URL)
}