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
	headerContentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
)

var uriBase = &url.URL{Host: "freebitco.in", Scheme: "https", Path: "/"}

const (
	urlRoot       = "https://freebitco.in"
	urlBase       = urlRoot + "/"
	urlSignUpPage = urlRoot + "/?op=signup_page"
	urlApi        = urlRoot + "/cgi-bin/api.pl"
	urlCaptcha    = "https://captchas.freebitco.in/botdetect/e/live/index.php"
)

var (
	regexSignupToken    = regexp.MustCompile(`signup_token\s=\s'(.*?)'`)
	regexBalance        = regexp.MustCompile(`<span id="balance">([0-9.]+)</span>`)
	regexDisableLottery = regexp.MustCompile(`id="disable_lottery_checkbox"\s(.*?)\s</div>`)
)
