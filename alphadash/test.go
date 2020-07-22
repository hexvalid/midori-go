package main

import (
	"fmt"
	"net/http"
)
import "github.com/gin-gonic/gin"

func main() {

	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	//router.LoadHTMLGlob("alphadash/resources/*")
	r.LoadHTMLFiles(
		"alphadash/resources/login.tmpl",
		"alphadash/resources/libs_head.tmpl",
		"alphadash/resources/libs_footer.tmpl")

	r.Static("/static", "alphadash/resources/static")
	r.StaticFile("/favicon.ico", "alphadash/resources/static/img/favicon.ico")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"error": false})
	})

	r.POST("/login", func(c *gin.Context) {
		fmt.Println(c.Params.ByName("username"))
		fmt.Println(c.Params.ByName("password"))

		if c.Param("username") == "root" && c.Param("password") == "1" {
			c.JSON(http.StatusOK, "ok")
			fmt.Println("OEKEOEKE")
		} else {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"error": true})
		}
	})
	r.Run(":8080")
}
