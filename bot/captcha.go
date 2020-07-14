package bot

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var captchaDownloaderClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	Timeout: 8 * time.Second,
}

func (a *Account) generateCaptcha() (string, error) {
	form := url.Values{}
	form.Add("op", "generate_captchasnet")
	form.Add("f", a.Browser.Fingerprint)
	form.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	req := a.newRequest(methodPost, urlApi, strings.NewReader(form.Encode()), true, urlBase)
	req.Header.Add(headerContentType, headerContentTypeFormUrlEncoded)
	res, err := a.execRequest(req)
	return strings.TrimSpace(res), err
}

func (a *Account) solveCaptcha() {

	random, _ := a.generateCaptcha()

	form := url.Values{}
	form.Add("random", random)
	req, _ := http.NewRequest(http.MethodPost, urlCaptcha, strings.NewReader(form.Encode()))
	req.Header.Add(headerAccept, browserAcceptImage)
	req.Header.Add("Referer", urlBase)
	req.Header.Add(headerContentType, headerContentTypeFormUrlEncoded)

	x, _ := captchaDownloaderClient.Do(req)
	defer x.Body.Close()
	body, _ := ioutil.ReadAll(x.Body)

	err := ioutil.WriteFile("/home/hexvalid/c.jpeg", body, 0644)
	if err != nil {
		// handle error
	}
}
