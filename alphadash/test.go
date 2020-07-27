package main

import (
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	gin.SetMode(gin.DebugMode)
	recaptcha.Init("6LcIqLQZAAAAAA_RFrQJjJy0-x5Oe4mbct9bHXZJ")

	r := gin.Default()
	store := cookie.NewStore([]byte("Y2FuxLFtIHlhxJ9tdXJ1bSDDp29rIHNldml5b3J1bSBzZW5p"))
	r.Use(sessions.Sessions("midori-auth", store))
	r.LoadHTMLFiles(
		"alphadash/resources/login.tmpl",
		"alphadash/resources/dashboard_base.tmpl",
		"alphadash/resources/libs_head.tmpl",
		"alphadash/resources/libs_footer.tmpl")

	r.Static("/static", "alphadash/resources/static")
	r.StaticFile("/favicon.ico", "alphadash/resources/static/img/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		if !checkValidLogin(c) {
			c.Redirect(http.StatusFound, "/login")
		} else {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	})
	//http://localhost:8080/static/css/AdminLTE.css
	//http://localhost:8080/dashboard/static/css/AdminLTE.css
	r.GET("/login", func(c *gin.Context) {
		if !checkValidLogin(c) {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"error": false})
		} else {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	})

	r.POST("/login", func(c *gin.Context) {
		cr, err := recaptcha.Confirm(c.ClientIP(), c.PostForm("g-recaptcha-response"))
		if err != nil || !cr {
			c.HTML(http.StatusUnauthorized, "login.tmpl",
				gin.H{"error": true, "errorMessage": "Captcha is incorrect"})
		} else if !(c.PostForm("username") == "root" && c.PostForm("password") == "toor") {
			c.HTML(http.StatusUnauthorized, "login.tmpl",
				gin.H{"error": true, "errorMessage": "Incorrect login credentials"})
		} else {
			session := sessions.Default(c)
			session.Options(sessions.Options{
				MaxAge: 3600 * 1, //1 go
			})
			session.Set("logged_user", c.PostForm("username"))
			session.Save()
			c.Redirect(http.StatusFound, "/dashboard")
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/login")
	})

	dashboard := r.Group("/dashboard")

	dashboard.Use(authRequired())
	{
		dashboard.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "dashboard_base.tmpl", gin.H{
				"title":    "Dash",
				"subtitle": "Main dash page template",
				"menu1":    true,
				"version":  "dev-preview",
			})
		})
	}

	r.Run(":8080")
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkValidLogin(c) {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

func checkValidLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	//todo: check username avaible
	return session.Get("logged_user") != nil
}
