package bot

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

var captchaDownloaderClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	Timeout: 8 * time.Second,
}

func (a *Account) generateCaptcha() (string, error) {
	req := a.newRequest(methodGet, urlApi, nil, true, urlBase)
	q := req.URL.Query()
	q.Add("op", "generate_captchasnet")
	q.Add("f", a.Browser.Fingerprint)
	q.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	req.URL.RawQuery = q.Encode()
	res, err := a.execRequest(req)
	return strings.TrimSpace(res), err
}

func (a *Account) downloadCaptcha(random string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, urlCaptcha, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("random", random)
	req.URL.RawQuery = q.Encode()
	req.Header.Add(headerAccept, browserAcceptImage)
	req.Header.Set(headerUserAgent, noUserAgent)
	x, err := captchaDownloaderClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer x.Body.Close()
	body, err := ioutil.ReadAll(x.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func (a *Account) isTrainedCaptcha(b []byte) (bool, error) {
	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		return false, err
	}
	whitePixelCount := 0
	grayPixelCount := 0
	for x := 0; x < img.Bounds().Size().X; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			color := img.At(x, y)
			r, g, b, _ := color.RGBA()
			if r == g && r == b {
				if r == 65535 {
					whitePixelCount++
				} else if r < 60909 {
					grayPixelCount++
				}
			}
		}
	}
	if (whitePixelCount > 7000 && whitePixelCount < 16762) && (grayPixelCount > 2120 && grayPixelCount < 4500) {
		return true, nil
	} else {
		return false, nil
	}
}

func (a *Account) solveCaptcha() (random, response string, err error) {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Looking up trained captcha...")
	var img []byte
	var found bool
	var wg sync.WaitGroup
	wg.Add(parallelCaptcha)
	for i := 0; i < parallelCaptcha; i++ {
		go func() {
			for !found {
				random_, err := a.generateCaptcha()
				if !found && err == nil {
					img_, err := a.downloadCaptcha(random_)
					if !found && err == nil {
						isTrained, err := a.isTrainedCaptcha(img_)
						if isTrained && !found && err == nil {
							found = true
							random = random_
							img = img_
							break
						}
					}
				}
			}
			wg.Done()
		}()
		time.Sleep(parallelCaptchaDelay)
	}
	wg.Wait()

	log.SInfo(fmt.Sprintf("%08d", a.ID), "Solving captcha via AI...")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "image.jpg")
	if err != nil {
		return "", "", err
	}
	if _, err = io.Copy(part, bytes.NewBuffer(img)); err != nil {
		return "", "", err
	}
	if err = writer.Close(); err != nil {
		return "", "", err
	}
	req, _ := http.NewRequest(http.MethodPost, urlAIServerSolve, body)
	req.Header.Add(headerContentType, writer.FormDataContentType())
	req.Header.Set(headerUserAgent, noUserAgent)
	res, err := captchaDownloaderClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	serverResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}
	response = strings.TrimSpace(string(serverResponse))
	return
}
