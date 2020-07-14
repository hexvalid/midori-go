package bot

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (a *Account) OpenBrowser() {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Opening Browser...")
	a.Browser.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar:     a.Browser.Jar,
		Timeout: browserTimeout,
	}
}

func (a *Account) newRequest(method, url string, body io.Reader, xreq bool, referer string) (req *http.Request) {
	req, _ = http.NewRequest(method, url, body)

	if !xreq {
		req.Header.Add(headerAccept, browserAcceptDefault)
	} else {
		req.Header.Add(headerAccept, browserAcceptAll)
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		req.Header.Add("Origin", urlBase)
		req.Header.Add("X-Csrf-Token", a.getCookieValue(cookieCsrfToken))
	}
	if referer != noReferer {
		req.Header.Add("Referer", referer)
	}

	req.Header.Add("Accept-Encoding", browserAcceptEncodingDefault)
	req.Header.Add("Accept-Language", a.Browser.AcceptLanguage)
	req.Header.Add("User-Agent", a.Browser.UserAgent)
	return
}

func (a *Account) execRequest(req *http.Request) (res string, err error) {
	rawRes, err := a.Browser.Client.Do(req)
	if err != nil {
		return
	}
	var body []byte
	defer rawRes.Body.Close()
	switch rawRes.Header.Get("Content-Encoding") {
	case "gzip":
		var gz *gzip.Reader
		if gz, err = gzip.NewReader(rawRes.Body); err != nil {
			return
		}
		defer gz.Close()
		body, err = ioutil.ReadAll(gz)
	case "deflate":
		fz := flate.NewReader(rawRes.Body)
		defer fz.Close()
		body, err = ioutil.ReadAll(fz)
	default:
		body, err = ioutil.ReadAll(rawRes.Body)
	}
	res = string(body)
	return
}

func (a *Account) getCookieValue(cookieName string) string {
	for _, cookie := range a.Browser.Jar.Cookies(uriBase) {
		if cookie.Name == cookieName {
			return cookie.Value
		}
	}
	return ""
}
