package bot

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/logger"
	"net/url"
	"strconv"
	"strings"
)

var log = logger.NewLog("Bot", color.FgBlue)

func (a *Account) Login(signUp bool) error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Loading initial page for login...")
	req := a.newRequest(methodGet, urlSignUpPage, nil, false, noReferer)
	body, _ := a.execRequest(req)

	a.solveCaptcha()

	form := url.Values{}
	form.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	if signUp {
		log.SInfo(fmt.Sprintf("%08d", a.ID), "Signing up...")
		form.Add("op", "signup_new")
		form.Add("email", a.Email)
		form.Add("fingerprint", a.Browser.Fingerprint)
		form.Add("referrer", strconv.Itoa(a.ReferrerID))
		form.Add("token", regexSignupToken.FindStringSubmatch(body)[1])
	} else {
		log.SInfo(fmt.Sprintf("%08d", a.ID), "Logging in...")
		form.Add("op", "login_new")
		form.Add("btc_address", a.Email)
		if a.Settings.EnableTFA && len(a.Settings.TFASecret) == 16 {
			tfaCode, err := GenerateTFA(a.Settings.TFASecret)
			if err != nil {
				return err
			}
			form.Add("tfa_code", tfaCode)
		} else {
			form.Add("tfa_code", "")
		}
	}
	form.Add("password", a.Password)

	req = a.newRequest(methodPost, urlBase, strings.NewReader(form.Encode()), true, urlSignUpPage)
	req.Header.Add(headerContentType, headerContentTypeFormUrlEncoded)

	return nil
}
