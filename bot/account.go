package bot

import (
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
	LastV3Date   time.Time
	Stats        string
	Settings     Settings
	Browser      Browser
	ActiveBoosts []Boost
	ReferrerID   int
	LoginTime    time.Time
	SignUpTime   time.Time
	SignUpIP     string
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
