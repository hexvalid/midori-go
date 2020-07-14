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
)

var uriBase = &url.URL{Host: "freebitco.in", Scheme: "https", Path: "/"}

const (
	urlRoot       = "https://freebitco.in"
	urlSubBase    = ".freebitco.in"
	urlBase       = urlRoot + "/"
	urlSignUpPage = urlRoot + "/?op=signup_page"
	urlApi        = urlRoot + "/cgi-bin/api.pl"
	urlCaptcha    = "https://captchas.freebitco.in/botdetect/e/live/index.php"
)

var (
	regexSignupToken  = regexp.MustCompile(`signup_token\s=\s'(.*?)'`)
	regexUserID       = regexp.MustCompile(`var userid = ([0-9.]+);`)
	regexBalance      = regexp.MustCompile(`<span id="balance">([0-9.]+)</span>`)
	regexRewardPoints = regexp.MustCompile(`user_reward_points\D+>([0-9,]+)</div>`)

	regexDisableLottery = regexp.MustCompile(`id="disable_lottery_checkbox"\s(.*?)\s</div>`)
)
