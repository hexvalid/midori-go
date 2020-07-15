package bot

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/getnada"
	"github.com/hexvalid/midori-go/logger"
	"github.com/hexvalid/midori-go/utils"
	"net/url"
	"strconv"
	"strings"
	"sync"
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
		captchaRandom, captchaResponse, err := a.solveCaptcha()
		if err != nil {
			return err
		}
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

	if strings.Contains(res, "invalid email address attached") {
		log.SInfo(fmt.Sprintf("%08d", a.ID), "%s %s", color.YellowString("(Warning)"),
			color.RedString("Invalid email address")+" error in homepage! Might cause a problem!")
		//todo: check email verification subsystem and do required thing! to best!
	}

	a.Balance, _ = strconv.ParseFloat(regexBalance.FindStringSubmatch(res)[1], 64)
	a.RewardPoints, _ = strconv.Atoi(regexRewardPoints.FindStringSubmatch(res)[1])
	a.Settings.SocketPassword = regexSocketPassword.FindStringSubmatch(res)[1]
	a.fpData.tokenName = regexTokenName.FindStringSubmatch(res)[1]
	a.fpData.secretTokenName = regexSecretTokenName.FindStringSubmatch(res)[1]
	a.fpData.secretToken = regexSecretToken.FindStringSubmatch(res)[1]
	a.fpData.token = regexToken.FindStringSubmatch(res)[1]
	a.fpData.fpAvailable, _ = strconv.Atoi(regexFPAvailable.FindStringSubmatch(res)[1])
	a.fpData.captchaType, _ = strconv.Atoi(regexCaptchaType.FindStringSubmatch(res)[1])

	if a.fpData.fpAvailable == 0 {
		fpTimeRemaining, _ := strconv.Atoi(regexFPTimeRemaining.FindStringSubmatch(res)[1])
		a.LastFPDate = time.Now().Add(-(time.Duration(fpTimeRemaining) * time.Second))
	}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		if err := a.loadStats("user_stats_initial", &a.Stats.UserStatsInitial); err != nil {
			log.SInfo(fmt.Sprintf("%08d", a.ID), "%s %s: %s", color.RedString("(Error)"),
				"Can't load user_stats_initial", err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		if err := a.loadStats("user_stats", &a.Stats.UserStats); err != nil {
			log.SInfo(fmt.Sprintf("%08d", a.ID), "%s %s: %s", color.RedString("(Error)"),
				"Can't load user_stats", err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		if err := a.recordFingerPrint(); err != nil {
			log.SInfo(fmt.Sprintf("%08d", a.ID), "%s %s: %s", color.RedString("(Error)"),
				"Can't record fingerprint", err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		if err := a.recordTimeOffset(); err != nil {
			log.SInfo(fmt.Sprintf("%08d", a.ID), "%s %s: %s", color.RedString("(Error)"),
				"Can't record time offset", err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		if err := a.getFPToken(); err != nil {
			log.SInfo(fmt.Sprintf("%08d", a.ID), "%s %s: %s", color.RedString("(Error)"),
				"Can't get FP token", err.Error())
		}
	}()

	wg.Wait()

	fmt.Println(a.fpData.fpAvailable, a.fpData.captchaType)
	return nil
}

func (a *Account) Roll() error {
	form := url.Values{}

	if a.fpData.captchaType == 1 {
		log.SInfo(fmt.Sprintf("%08d", a.ID), "Rolling without captcha...")
		form.Add("pwc", "1")
	} else if a.fpData.captchaType == 11 {
		log.SInfo(fmt.Sprintf("%08d", a.ID), "Rolling with captcha...")
		form.Add("pwc", "0")
		form.Add("g_recaptcha_response", "")
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			captchaRandom, captchaResponse, _ := a.solveCaptcha()
			form.Add("botdetect_random", captchaRandom)
			form.Add("botdetect_response", captchaResponse)
		}()

		go func() {
			defer wg.Done()
			captchaRandom, captchaResponse, _ := a.solveCaptcha()
			form.Add("botdetect_random2", captchaRandom)
			form.Add("botdetect_response2", captchaResponse)
		}()

		wg.Wait()
	}

	form.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	form.Add("op", "free_play")
	form.Add("fingerprint", a.Browser.Fingerprint)
	form.Add("client_seed", utils.RandomStringInRunes(16, utils.LetterBytes))
	form.Add("fingerprint2", strconv.Itoa(a.Browser.Fingerprint2))
	form.Add(a.fpData.tokenName, a.fpData.token)
	form.Add(a.fpData.secretTokenName, fmt.Sprintf("%x", sha256.Sum256(([]byte(a.fpData.fpToken))[:])))
	req := a.newRequest(methodPost, urlBase, strings.NewReader(form.Encode()), true, urlHomePage)
	req.Header.Add(headerContentType, headerContentTypeFormUrlEncoded)

	fmt.Println(a.execRequest(req))

	return nil
}

func (a *Account) loadStats(statType string, v interface{}) error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Loading stats: %s...", color.YellowString(statType))
	req := a.newRequest(methodGet, urlStatsNewPrivate, nil, true, urlBase)
	q := req.URL.Query()
	q.Add("u", strconv.Itoa(a.ID))
	q.Add("p", a.Settings.SocketPassword)
	q.Add("f", statType)
	q.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	req.URL.RawQuery = q.Encode()
	res, err := a.execRequest(req)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(res), &v); err != nil {
		return err
	}
	return nil
}

func (a *Account) recordFingerPrint() error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Recording fingerprint...")
	req := a.newRequest(methodGet, urlApi, nil, true, urlBase)
	q := req.URL.Query()
	q.Add("op", "record_fingerprint")
	q.Add("fingerprint", a.Browser.Fingerprint)
	q.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	req.URL.RawQuery = q.Encode()
	res, err := a.execRequest(req)
	if err != nil {
		return err
	} else if res != "1" {
		return fmt.Errorf("unknown record fingerprint response: %s", res)
	}
	return nil
}

func (a *Account) recordTimeOffset() error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Recording time offset...")
	req := a.newRequest(methodGet, urlApi, nil, true, urlBase)
	q := req.URL.Query()
	q.Add("op", "record_user_data")
	q.Add("type", "time_offset")
	q.Add("value", strconv.Itoa(a.Browser.TimeOffset))
	q.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	req.URL.RawQuery = q.Encode()
	_, err := a.execRequest(req)
	return err
}

func (a *Account) getFPToken() error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Getting FP token...")
	req := a.newRequest(methodGet, urlFpCheck, nil, true, urlBase)
	q := req.URL.Query()
	q.Add("s", a.fpData.secretTokenName)
	q.Add("csrf_token", a.getCookieValue(cookieCsrfToken))
	req.URL.RawQuery = q.Encode()
	res, err := a.execRequest(req)
	if err != nil {
		return err
	} else if len(res) < 12 {
		return fmt.Errorf("unknown fp token: %s", res)
	}
	a.fpData.fpToken = res
	return nil
}

func (a *Account) VerifyEmail() error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Requesting verification email...")
	req := a.newRequest(methodGet, urlMarkEmailValid, nil, false, urlBase)
	res, err := a.execRequest(req)
	if err != nil {
		return err
	}
	if !strings.Contains(res, "email sent") {
		if strings.Contains(res, "every 24 hours") {
			return errors.New("verification mail already sended in 24 hours")
		} else {
			return fmt.Errorf("unknown response during request verification mail: %s", res)
		}
	} else {
		a.Settings.EmailVerificationRequestTime = time.Now()
	}
	time.Sleep(emailVerificationWaitInitial)
	for i := 0; i < emailVerificationMaxCheckCount; i++ {
		inbox, err := getnada.GetInbox(a.Email)
		if err != nil {
			return err
		}
		for _, mail := range inbox {
			if strings.Contains(mail.FromEmail, "@freebitco.in") &&
				strings.Contains(strings.ToLower(mail.Subject), "verification") {
				if err = mail.Load(); err != nil {
					return err
				}
				a.Settings.EmailVerificationLink1 = regexMailImage1.FindStringSubmatch(mail.HTML)[1]
				a.Settings.EmailVerificationLink2 = regexMailImage2.FindStringSubmatch(mail.HTML)[1]
				if err = a.VerifyEmailLinks(); err != nil {
					return err
				}
				a.Settings.EmailVerified = true
				mail.Delete()
				return nil
			}
		}
		time.Sleep(emailVerificationWaitInterval)
	}
	return errors.New("verification timed out")
}

func (a *Account) VerifyEmailLinks() error {
	log.SInfo(fmt.Sprintf("%08d", a.ID), "Verifying email...")
	dummyReq := a.newRequest(methodGet, a.Settings.EmailVerificationLink1, nil, false, noReferer)
	dummyReq.Header.Add(headerAccept, browserAcceptImage)
	if _, err := a.execRequest(dummyReq); err != nil {
		return err
	}
	dummyReq = a.newRequest(methodGet, a.Settings.EmailVerificationLink2, nil, false, noReferer)
	dummyReq.Header.Add(headerAccept, browserAcceptImage)
	if _, err := a.execRequest(dummyReq); err != nil {
		return err
	}
	a.Settings.EmailVerificationTime = time.Now()
	return nil
}
