package bot

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/logger"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var log = logger.NewLog("Bot", color.FgBlue)

func (a *Account) Login(signUp bool) error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Loading initial page for login...")
	req := a.newRequest(methodGet, urlSignUpPage, nil, false, noReferer)
	res, err := a.execRequest(req)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	if signUp {
		captchaRandom, captchaResponse, _ := a.solveCaptcha()
		log.SInfo(fmt.Sprintf("%08d", a.ID), "Signing up...")
		form.Add("op", "signup_new")
		form.Add("email", a.Email)
		form.Add("fingerprint", a.Browser.Fingerprint)
		form.Add("referrer", strconv.Itoa(a.ReferrerID))
		form.Add("token", regexSignupToken.FindStringSubmatch(res)[1])
		form.Add("botdetect_random", captchaRandom)
		form.Add("botdetect_response", captchaResponse)
		a.SignUpTime = time.Now()
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

	res, err = a.execRequest(req)
	if err != nil {
		return err
	}

	body := strings.Split(res, ":")

	if body[0] == "s" {
		a.LoginTime = time.Now()
		a.BTCAddress = body[1]
		a.ID, _ = strconv.Atoi(body[3])
		a.addCookie("password", body[2])
		a.addCookie("btc_address", a.BTCAddress)
		a.addCookie("have_account", "1")
		log.SInfo(fmt.Sprintf("%08d", a.ID), color.GreenString("Successfully logged in."))
	} else if body[0] == "e" {
		return errors.New(strings.ToLower(body[1]))
	} else {
		return fmt.Errorf("unknown response: %s", res)
	}

	return nil
}

func (a *Account) Home() error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Loading home...")
	req := a.newRequest(methodGet, urlBase, nil, false, urlBase)
	res, _ := a.execRequest(req)

	if strconv.Itoa(a.ID) != regexUserID.FindStringSubmatch(res)[1] {
		return errors.New("unable verify user id at home")
	}

	a.Balance, _ = strconv.ParseFloat(regexBalance.FindStringSubmatch(res)[1], 64)
	a.RewardPoints, _ = strconv.Atoi(regexRewardPoints.FindStringSubmatch(res)[1])

	fmt.Println(res)
	fmt.Println(a.Balance)
	fmt.Println(a.RewardPoints)

	return nil
}
