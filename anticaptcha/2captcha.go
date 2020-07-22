package anticaptcha

// package twocaptcha provides a Golang client for https://2captcha.com/

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ApiURL is the url of the 2captcha API endpoint
var ApiURL = "https://2captcha.com/in.php"

// ResultURL is the url of the 2captcha result API endpoint
var ResultURL = "https://2captcha.com/res.php"

const ApiKey = "8200c6ea53180c85a39002b2b37e4d03"

var Client = http.Client{}

func SolveRecaptchaV3(siteURL, recaptchaKey, action string) (string, error) {
	captchaId, err := apiRequest(
		ApiURL,
		map[string]string{
			"googlekey": recaptchaKey,
			"pageurl":   siteURL,
			"method":    "userrecaptcha",
			"version":   "v3",
			"action":    action,
			"min_score": "0.9",
		},
		0,
		3,
	)

	if err != nil {
		return "", err
	}

	return apiRequest(
		ResultURL,
		map[string]string{
			"googlekey": recaptchaKey,
			"pageurl":   siteURL,
			"method":    "userrecaptcha",
			"id":        captchaId,
			"action":    "get",
		},
		5,
		20,
	)
}

func apiRequest(URL string, params map[string]string, delay time.Duration, retries int) (string, error) {
	if retries <= 0 {
		return "", errors.New("Maximum retries exceeded")
	}
	time.Sleep(delay * time.Second)
	form := url.Values{}
	form.Add("key", ApiKey)
	for k, v := range params {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if strings.Contains(string(body), "CAPCHA_NOT_READY") {
		return apiRequest(URL, params, delay, retries-1)
	}
	if !strings.Contains(string(body), "OK|") {
		return "", errors.New("Invalid respponse from 2captcha: " + string(body))
	}
	return string(body[3:]), nil
}
