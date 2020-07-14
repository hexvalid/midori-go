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
	RollCount    int
	LastFPDate   time.Time
	Stats        string
	Settings     Settings
	Browser      Browser
	ActiveBoosts []Boost
	ReferrerID   int
	LoginTime    time.Time
	SignUpTime   time.Time
	Proxy        string
}

type Settings struct {
	EmailVerified   bool
	DisableLottery  bool
	DisableInterest bool
	EnableTFA       bool
	TFASecret       string
	EmailSubs       string
}

type Browser struct {
	AcceptLanguage string
	UserAgent      string
	Fingerprint    string
	Fingerprint2   int
	Jar            *cookiejar.Jar
	Client         *http.Client
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
	a.Browser.Jar, _ = cookiejar.New(nil)
	return
}
