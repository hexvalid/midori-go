package getnada

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/logger"
	"github.com/hexvalid/midori-go/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	apiBaseUrl  = "https://getnada.com/api/v1"
	apiDomains  = "/domains"
	apiInboxes  = "/inboxes/"
	apiMessages = "/messages/"
)

var (
	log    = logger.NewLog("GetNada", color.FgHiYellow)
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		},
		Timeout: 10 * time.Second,
	}
	domains []string
)

func getDomains() (err error) {
	log.Info("Getting domains...")
	var res *http.Response
	var body []byte
	var result []resDomains
	domains = []string{}
	if res, err = client.Get(apiBaseUrl + apiDomains); err != nil {
		return err
	}
	defer res.Body.Close()
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return err
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}
	for _, i := range result {
		domains = append(domains, i.Name)
	}
	log.Info("%s domain received.", color.YellowString(strconv.Itoa(len(domains))))
	return
}

func GenerateMail() (string, error) {
	if len(domains) < 1 {
		if err := getDomains(); err != nil {
			return "", err
		}
	}
	rand.Seed(time.Now().Unix())
	address := fmt.Sprintf("%s@%s", utils.GenerateUsername(), domains[rand.Intn(len(domains))])
	log.Info("Mail address generated: %s", color.YellowString(address))
	return address, nil
}

func GetInbox(address string) (mails []Mail, err error) {
	log.Info("Getting inbox of %s...", color.YellowString(address))
	var res *http.Response
	var body []byte
	var result resInboxes
	domains = []string{}
	if res, err = client.Get(apiBaseUrl + apiInboxes + address); err != nil {
		return
	}
	defer res.Body.Close()
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	mails = result.Msgs
	log.Info("%s inbox received.", color.YellowString(strconv.Itoa(len(mails))))
	return
}

func (m *Mail) Load() (html string, err error) {
	log.SInfo(fmt.Sprintf(m.UID[:6]), "Loading mail...")
	var res *http.Response
	var body []byte
	if res, err = client.Get(apiBaseUrl + apiMessages + m.UID); err != nil {
		return
	}
	defer res.Body.Close()
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, m); err != nil {
		return
	}
	return m.HTML, nil
}

func (m *Mail) Delete() (err error) {
	log.SInfo(fmt.Sprintf(m.UID[:6]), "Deleting mail...")
	var req *http.Request
	var res *http.Response
	if req, err = http.NewRequest("DELETE", apiBaseUrl+apiMessages+m.UID, nil); err != nil {
		return
	}
	if res, err = client.Do(req); err != nil {
		return
	}
	res.Body.Close()
	return nil
}
