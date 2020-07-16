package bot

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/tormdr"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ProxyTypeSocks5 = "socks5"
	ProxyTypeTor    = "tor"
)

type ProxyConfig struct {
	Enabled      bool    `json:"enabled"`
	CurrentProxy Proxy   `json:"currentProxy,omitempty"`
	ProxyHistory []Proxy `json:"proxyHistory,omitempty"`
}

type Proxy struct {
	Type     string `json:"type,omitempty"`
	Address  string `json:"address,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (a *Account) checkIP() (ip string, latency int, err error) {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Testing IP Address...")
	var req *http.Request
	var res *http.Response
	var body []byte

	if req, err = http.NewRequest(http.MethodGet, ipCheckServer, nil); err != nil {
		return
	}
	req.Header.Set("User-Agent", "")
	start := time.Now()
	if res, err = a.client.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	latency = int(time.Since(start) / time.Millisecond)
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	ip = strings.TrimSpace(string(body))
	log.SInfo(fmt.Sprintf("%08d", a.ID), "IP Address tested. IP: %s, Latency: %s ms.",
		color.YellowString(ip), color.YellowString(strconv.Itoa(latency)))
	return
}

func (a *Account) PlugNewTorAddress(excludedIPs []string) (err error) {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Plugging new Tor IP...")
	ip, err := tormdr.FindExitNode(excludedIPs, 8000, true, true, true)
	if err != nil {
		return err
	} else if len(ip) < 8 {
		return errors.New("non valid ip address")
	}
	a.Proxy.Enabled = true
	a.Proxy.CurrentProxy.Type = ProxyTypeTor
	a.Proxy.CurrentProxy.Address = ip
	return nil
}
