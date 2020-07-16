package bot

import (
	"net/url"
	"regexp"
	"time"
)

const (
	browserTimeout                  = 30 * time.Second
	browserAcceptDefault            = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	browserAcceptAll                = "*/*"
	browserAcceptImage              = "image/webp,*/*"
	browserAcceptEncodingDefault    = "gzip, deflate"
	methodGet                       = "GET"
	methodPost                      = "POST"
	noReferer                       = ""
	cookieCsrfToken                 = "csrf_token"
	headerAccept                    = "Accept"
	headerContentType               = "Content-Type"
	headerUserAgent                 = "User-Agent"
	noUserAgent                     = ""
	headerContentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
	parallelCaptcha                 = 4
	parallelCaptchaDelay            = 100 * time.Millisecond
	urlAIServerSolve                = "http://40.70.243.118:4200/solve/"
	emailVerificationMaxCheckCount  = 6
	emailVerificationWaitInterval   = 10 * time.Second
	emailVerificationWaitInitial    = 5 * time.Second
	recaptchaV3Key                  = "6Lc1kXIUAAAAAPP7OeuycKWZ-t4br4Rh3XvqWUGd"
)

var uriBase = &url.URL{Host: "freebitco.in", Scheme: "https", Path: "/"}

const (
	urlRoot            = "https://freebitco.in"
	urlSubBase         = ".freebitco.in"
	urlBase            = urlRoot + "/"
	urlHomePage        = urlRoot + "/?op=home"
	urlSignUpPage      = urlRoot + "/?op=signup_page"
	urlApi             = urlRoot + "/cgi-bin/api.pl"
	urlStatsNewPrivate = urlRoot + "/stats_new_private/"
	urlFpCheck         = urlRoot + "/cgi-bin/fp_check.pl"
	urlMarkEmailValid  = urlRoot + "/?op=mark_email_valid"
	urlCaptcha         = "https://captchas.freebitco.in/botdetect/e/live/index.php"
)

var (
	regexSignupToken     = regexp.MustCompile(`signup_token\s=\s'(.*?)'`)
	regexUserID          = regexp.MustCompile(`var userid = ([0-9.]+);`)
	regexBalance         = regexp.MustCompile(`<span id="balance">([0-9.]+)</span>`)
	regexRewardPoints    = regexp.MustCompile(`user_reward_points\D+>([0-9,]+)</div>`)
	regexSocketPassword  = regexp.MustCompile(`socket_password = '(.*?)'`)
	regexTokenName       = regexp.MustCompile(`token_name = '(.*?)'`)
	regexSecretTokenName = regexp.MustCompile(`tcGiQefA = '(.*?)'`)
	regexSecretToken     = regexp.MustCompile(`um2VHVjSZ = \"(.*?)\"`)
	regexToken           = regexp.MustCompile(`token1 = '(.*?)'`)
	regexFPAvailable     = regexp.MustCompile(`free_play = (.*?);`)
	regexCaptchaType     = regexp.MustCompile(`captcha_type = (.*?);`)
	regexFPTimeRemaining = regexp.MustCompile(`\.free_play_time_remaining'\).countdown\({until: ([0-9.+-]+),`)
	regexMailImage1      = regexp.MustCompile(`<img src='(.*?)'`)
	regexMailImage2      = regexp.MustCompile(`<img src=\"(.*?)\"`)

	regexDisableLottery = regexp.MustCompile(`id="disable_lottery_checkbox"\s(.*?)\s</div>`)
)
