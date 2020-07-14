package bot

import (
	"github.com/hexvalid/midori-go/getnada"
	"github.com/hexvalid/midori-go/utils"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Account struct {
	ID           int
	Email        string
	BTCAddress   string
	Password     string
	Balance      float64
	RewardPoints int
	FPCount      int
	LastFPDate   time.Time
	Stats        string
	Settings     Settings
	Browser      Browser
	ActiveBoosts []Boost
	ReferrerID   int
	LoginTime    time.Time
	SignUpTime   time.Time
	Serial       string
	Proxy        string

	jar    *cookiejar.Jar
	client *http.Client
}

type Settings struct {
	EmailVerified   bool   `json:"emailVerified"`
	DisableLottery  bool   `json:"disableLottery"`
	DisableInterest bool   `json:"disableInterest"`
	EnableTFA       bool   `json:"enableTFA"`
	TFASecret       string `json:"tfaSecret"`
	EmailSubs       string `json:"emailSubs"`
}

type Browser struct {
	AcceptLanguage string `json:"acceptLanguage"`
	UserAgent      string `json:"userAgent"`
	Fingerprint    string `json:"fingerprint"`
	Fingerprint2   int    `json:"fingerprint2"`
}

func GenerateNewAccount(referrerID int) (a Account, err error) {
	log.Info("Generating new account...")
	if a.Email, err = getnada.GenerateMail(); err != nil {
		return
	}
	a.Password = utils.RandomStringInRunes(utils.RandomInt(8, 16), utils.LetterBytes)
	a.ReferrerID = referrerID
	a.Browser.AcceptLanguage = utils.RandomStringInArray(utils.AcceptLanguages)
	a.Browser.UserAgent = utils.RandomStringInArray(utils.UserAgents)
	a.Browser.Fingerprint = utils.RandomStringInRunes(32, utils.BaseBytes)
	a.Browser.Fingerprint2 = utils.RandomInt(1111111111, 9999999999)
	a.jar, _ = cookiejar.New(nil)
	return
}
