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
	Stats        Stats
	Settings     Settings
	Browser      Browser
	ActiveBoosts []Boost
	ReferrerID   int
	LoginTime    time.Time
	SignUpTime   time.Time
	Serial       string
	Proxy        ProxyConfig

	jar    *cookiejar.Jar
	client *http.Client
	fpData FPData
}

type Settings struct {
	DisableLottery               bool      `json:"disableLottery"`
	DisableInterest              bool      `json:"disableInterest"`
	EnableTFA                    bool      `json:"enableTFA"`
	TFASecret                    string    `json:"tfaSecret"`
	EmailSubs                    string    `json:"emailSubs"`
	SocketPassword               string    `json:"socketPassword"`
	EmailVerified                bool      `json:"emailVerified"`
	EmailVerificationRequestTime time.Time `json:"emailVerificationRequestTime"`
	EmailVerificationTime        time.Time `json:"emailVerificationTime"`
	EmailVerificationLink1       string    `json:"emailVerificationLink1"`
	EmailVerificationLink2       string    `json:"emailVerificationLink2"`
	RecordRecaptchaV3Count       int       `json:"recordRecaptchaV3Count"`
	RecordRecaptchaV3LastDate    time.Time `json:"recordRecaptchaV3LastDate"`
	RecordRecaptchaV3LastScore   float64   `json:"recordRecaptchaV3LastScore"`
}

type Browser struct {
	AcceptLanguage string `json:"acceptLanguage"`
	UserAgent      string `json:"userAgent"`
	Fingerprint    string `json:"fingerprint"`
	Fingerprint2   int    `json:"fingerprint2"`
	TimeOffset     int    `json:"timeOffset"`
}

type FPData struct {
	tokenName       string
	token           string
	secretTokenName string
	secretToken     string
	fpToken         string

	fpAvailable int
	captchaType int
}

type Stats struct {
	UserStatsInitial `json:"userStatsInitial"`
	UserStats        `json:"userStats"`
}

type UserStatsInitial struct {
	//PaymentsSent 	[]interface{} `json:"payments_sent"`
	//Deposits     	[]interface{} `json:"deposits"`
	User struct {
		ReferralCommissionsEarned string `json:"referral_commissions_earned"`
		GrossBalance              string `json:"gross_balance"`
		LotterySpent              string `json:"lottery_spent"`
		FreeSpinsPlayed           int    `json:"free_spins_played"`
		JackpotWinnings           string `json:"jackpot_winnings"`
		TotalPayouts              int    `json:"total_payouts"`
		JackpotSpent              string `json:"jackpot_spent"`
		PaidWinnings              string `json:"paid_winnings"`
		FreeWinnings              string `json:"free_winnings"`
		PaidSpinsPlayed           int    `json:"paid_spins_played"`
		//TotalDeposits           interface{} `json:"total_deposits"`
	} `json:"user"`
}

type UserStats struct {
	LotteryTickets string `json:"lottery_tickets"`
	UserExtras     struct {
		LotterySpent string `json:"lottery_spent"`
		//MultiplyCommissionsEarned int    `json:"multiply_commissions_earned,omitempty"`
		JackpotSpent string `json:"jackpot_spent"`
	} `json:"user_extras"`
	User struct {
		FreeWinnings int `json:"free_winnings"`
	} `json:"user"`
	UnblockGbr struct {
		LotteryToUnblock int    `json:"lottery_to_unblock"`
		WagerToUnblock   string `json:"wager_to_unblock"`
		JackpotToUnblock string `json:"jackpot_to_unblock"`
		DepositToUnblock string `json:"deposit_to_unblock"`
	} `json:"unblock_gbr"`
	WagerContest struct {
		WagerPersonal      int `json:"wager_personal"`
		RefContestPersonal int `json:"ref_contest_personal"`
	} `json:"wager_contest"`
	LamboContestRound int `json:"lambo_contest_round"`
	//InstantPaymentRequests []interface{} `json:"instant_payment_requests"`
	LamboLotteryTickets interface{} `json:"lambo_lottery_tickets"`
	NoCaptchaGbr        struct {
		LotteryToUnblock int    `json:"lottery_to_unblock"`
		WagerToUnblock   string `json:"wager_to_unblock"`
		JackpotToUnblock string `json:"jackpot_to_unblock"`
		DepositToUnblock string `json:"deposit_to_unblock"`
	} `json:"no_captcha_gbr"`
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
	a.Browser.TimeOffset = utils.RandomIntInArray(utils.TimeOffsets)
	a.jar, _ = cookiejar.New(nil)
	return
}
