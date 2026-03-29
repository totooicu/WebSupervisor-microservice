package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ParseHTML(content string, keys []string) []string {
	var results []string
	
	for _, key := range keys {
		parts := strings.Split(key, ",")
		if len(parts) >= 2 {
			left := parts[0]
			right := parts[1]
			values := GetMid(content, left, right, 0)
			results = append(results, values...)
		}
	}
	
	return results
}

func ParseJSON(content string, jsonKeys []string) []interface{} {
	var data interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil
	}
	
	var results []interface{}
	
	for _, key := range jsonKeys {
		value := GetJSONValue(data, strings.Split(key, "."))
		if value != nil {
			results = append(results, value)
		}
	}
	
	return results
}

func GetMid(content, left, right string, count int) []string {
	var results []string
	index := 0
	
	for i := 0; i < count || count == 0; i++ {
		leftIdx := strings.Index(content[index:], left)
		if leftIdx == -1 {
			break
		}
		leftIdx += index + len(left)
		
		rightIdx := strings.Index(content[leftIdx:], right)
		if rightIdx == -1 {
			break
		}
		rightIdx += leftIdx
		
		result := content[leftIdx:rightIdx]
		results = append(results, strings.TrimSpace(result))
		index = rightIdx + len(right)
	}
	
	return results
}

func GetJSONValue(data interface{}, path []string) interface{} {
	for _, key := range path {
		switch v := data.(type) {
		case map[string]interface{}:
			if val, ok := v[key]; ok {
				data = val
			} else {
				return nil
			}
		case []interface{}:
			if idx, err := fmt.Sscanf(key, "%d"); err == nil {
				if idx >= 0 && idx < len(v) {
					data = v[idx]
				} else {
					return nil
				}
			} else {
				return nil
			}
		default:
			return nil
		}
	}
	
	return data
}

func StringArrayMustCompileStringArray(source []string, patterns []string) []string {
	var results []string
	
	for _, src := range source {
		for _, pattern := range patterns {
			matched, err := regexp.MatchString(pattern, src)
			if err == nil && matched {
				results = append(results, src)
				break
			}
		}
	}
	
	return results
}