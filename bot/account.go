package bot

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Account struct {
	ID             int
	Email          string
	Password       string
	Balance        float64
	RewardPoints   int
	RollCount      int
	NextRollTime   time.Time
	EmailVerified  bool
	TFAEnabled     bool
	TFASecret      string
	BTCAddress     string
	DisableLottery bool
	UserAgent      string
	Fingerprint    string
	Fingerprint2   int
	LoginTime      time.Time
	SignUpTime     time.Time
	SignUpIP       string
	Referrer       int
	ActiveBoosts   []Boost
	Transfers      string
	Jar            *cookiejar.Jar
	Client         *http.Client
}
