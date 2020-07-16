package bot

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hexvalid/midori-go/tormdr"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

const ipCheckServer string = "http://checkip.amazonaws.com"

func (a *Account) OpenBrowser(tormdr *tormdr.TorMDR) error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Opening Browser...")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if a.Proxy.Enabled {
		if a.Proxy.CurrentProxy.Type == ProxyTypeTor {
			if tormdr == nil {
				return errors.New("tormdr is not attached")
			}
			if err := tormdr.SetExitNode(a.Proxy.CurrentProxy.Address); err != nil {
				return err
			}

			ip1, _, err := tormdr.CheckIP()
			if err != nil {
				return err
			}

			transport.Proxy = http.ProxyURL(tormdr.Proxy)

			//todo: dublicated
			a.client = &http.Client{
				Transport: transport,
				Jar:       a.jar,
				Timeout:   browserTimeout,
			}

			ip2, _, err := a.checkIP()
			if err != nil {
				return err
			}

			if !(a.Proxy.CurrentProxy.Address == ip2 && ip1 == ip2) {
				return errors.New("ip address mismatched")
			}

		}
	} else {
		a.client = &http.Client{
			Transport: transport,
			Jar:       a.jar,
			Timeout:   browserTimeout,
		}
	}

	return nil
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
	rawRes, err := a.client.Do(req)
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
	for _, cookie := range a.jar.Cookies(uriBase) {
		if cookie.Name == cookieName {
			return cookie.Value
		}
	}
	return ""
}

func (a *Account) addCookie(name, value string) {
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		Domain: urlSubBase,
	})
	a.jar.SetCookies(uriBase, cookies)
}

func (a *Account) JarToString() string {
	var buffer bytes.Buffer
	cookies := a.jar.Cookies(uriBase)
	for i := 0; i < len(cookies); i++ {
		buffer.WriteString(cookies[i].Name)
		buffer.WriteString("=")
		buffer.WriteString(cookies[i].Value)
		if i != len(cookies)-1 {
			buffer.WriteString("; ")
		}
	}
	return buffer.String()
}

func (a *Account) StringToJar(s string) {
	a.jar, _ = cookiejar.New(nil)
	var cookies []*http.Cookie
	ss := strings.Split(s, "; ")
	for i := range ss {
		sscc := strings.Split(ss[i], "=")
		cookie := &http.Cookie{
			Name:  sscc[0],
			Value: sscc[1],
		}
		cookies = append(cookies, cookie)
	}
	a.jar.SetCookies(uriBase, cookies)
}
