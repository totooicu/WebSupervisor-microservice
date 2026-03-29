package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpClient struct {
	url    string
	header map[string]string
}

func NewHttpClient(url string, header map[string]string) *HttpClient {
	return &HttpClient{
		url:    url,
		header: header,
	}
}

func (h *HttpClient) Get() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", h.url, nil)
	if err != nil {
		return "", err
	}

	for k, v := range h.header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (h *HttpClient) Post(body map[string]interface{}, strPayload string) (string, error) {
	client := &http.Client{}

	var jsonData []byte
	var err error

	if strPayload != "" {
		jsonData = []byte(strPayload)
	} else {
		jsonData, err = json.Marshal(body)
		if err != nil {
			return "", err
		}
	}

	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	for k, v := range h.header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
