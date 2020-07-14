package tormdr

import (
	"crypto/tls"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const ipCheckServer string = "http://checkip.amazonaws.com"

func (tormdr *TorMDR) CheckIP() (ip string, latency int, err error) {
	log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Testing Exit Node...")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(tormdr.proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	var req *http.Request
	var res *http.Response
	var body []byte

	if req, err = http.NewRequest(http.MethodGet, ipCheckServer, nil); err != nil {
		return
	}
	req.Header.Set("User-Agent", "")
	start := time.Now()
	if res, err = client.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	latency = int(time.Since(start) / time.Millisecond)
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	ip = strings.TrimSpace(string(body))
	log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Exit Node tested. IP: %s, Latency: %s ms.",
		color.YellowString(ip), color.YellowString(strconv.Itoa(latency)))
	return
}
