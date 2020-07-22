package anticaptcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/logger"
	"net/http"
	"net/url"
	"time"
)

var (
	baseURL       = &url.URL{Host: "api.anti-captcha.com", Scheme: "https", Path: "/"}
	checkInterval = 4 * time.Second
	timeout       = 60 * time.Second
	apiKey        = "da8b701dfebc84853a983fe9c7794a9d"
	minScore      = 0.9
	log           = logger.NewLog("AntiCaptcha", color.FgHiBlue)
)

func createTaskRecaptcha(websiteURL string, recaptchaKey, action string) (float64, error) {
	// Mount the data to be sent
	body := map[string]interface{}{
		"clientKey": apiKey,
		"task": map[string]interface{}{
			"type":       "RecaptchaV3TaskProxyless",
			"websiteURL": websiteURL,
			"websiteKey": recaptchaKey,
			"minScore":   minScore,
			"pageAction": action,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}

	// Make the request
	u := baseURL.ResolveReference(&url.URL{Path: "/createTask"})
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Decode response
	responseBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseBody)
	// TODO treat api errors and handle them properly
	if _, ok := responseBody["taskId"]; ok {
		if taskId, ok := responseBody["taskId"].(float64); ok {
			return taskId, nil
		}

		return 0, errors.New("task number of irregular format")
	}

	return 0, errors.New("task number not found in server response")
}

func getTaskResult(taskID float64) (map[string]interface{}, error) {
	// Mount the data to be sent
	body := map[string]interface{}{
		"clientKey": apiKey,
		"taskId":    taskID,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Make the request
	u := baseURL.ResolveReference(&url.URL{Path: "/getTaskResult"})
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	responseBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseBody)
	return responseBody, nil
}

func SendRecaptcha(websiteURL, recaptchaKey, action string) (string, error) {
	log.Info("Creating recaptcha task...")
	taskID, err := createTaskRecaptcha(websiteURL, recaptchaKey, action)
	if err != nil {
		return "", err
	}

	check := time.NewTicker(10 * time.Second)
	timeout := time.NewTimer(timeout)

	for {
		select {
		case <-check.C:
			response, err := getTaskResult(taskID)
			if err != nil {
				return "", err
			}
			if response["status"] == "ready" {
				return response["solution"].(map[string]interface{})["gRecaptchaResponse"].(string), nil
			}
			check = time.NewTicker(checkInterval)
		case <-timeout.C:
			return "", errors.New("antiCaptcha check result timeout")
		}
	}
}
