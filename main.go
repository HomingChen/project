package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var authCode string = ""

func main() {
	loadEnvFile()

	router := gin.Default()

	store := cookie.NewStore([]byte(getEnvValue("sessionKey")))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   getEnvValue("sessionDomain"),
		MaxAge:   60,
		Secure:   false,
		HttpOnly: false,
	})
	router.Use(sessions.Sessions("status", store))

	router.LoadHTMLGlob("template/*")
	router.Static("/static", "./static")
	router.GET("/", indexHandler)
	router.GET("/login", loginHandler)
	router.GET("/logout", logoutHandler)
	router.GET("/api/getAuthURL", authURLHandler)
	router.GET("/api/getGoogleAuthCode", getGoogleAuthHandler)
	router.GET("/api/getNotionAuthCode", getNotionAuthHandler)
	router.GET("/api/getMicrosoftAuthCode", getMicrosoftAuthHandler)
	router.Run(":8080")
}

func indexHandler(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Printf("[NOTE] session: %v\n", session.Get("id"))
	c.HTML(http.StatusOK, "index.html", "")
}

func loginHandler(c *gin.Context) {
	session := sessions.Default(c)
	checkLogin := session.Get("id")
	fmt.Printf("[NOTE] session: %T, %v\n", checkLogin, checkLogin)
	if checkLogin == "logined" {
		c.Redirect(302, "/")
	} else {
		c.HTML(http.StatusOK, "login.html", "")
	}
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	fmt.Printf("[NOTE] session: %v\n", session.Get("status"))
	c.Redirect(302, "/")
}

func authURLHandler(c *gin.Context) {
	authType := c.Query("authType")
	var authURL string
	switch authType {
	case "google":
		config := formGoogleAuthConfig()
		authURL = getGoogleAuthURL(config)
	case "notion":
		authURL = getEnvValue("notionAuthURL")
	case "microsoft":
		authURL = "/api/getMicrosoftAuthCode"
	default:
		authURL = "no authType"
	}
	// c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	c.JSON(200, authURL)
}

func getNotionAuthHandler(c *gin.Context) {
	code := c.Query("code")
	fmt.Printf("NotionCode: %s\n", code)
	client := requestNotionAuth(code)
	// showNotionUserInfo(client)
	userInfo := getNotionUserInfo(client)
	addLoginCookie(userInfo, c)

	// session := sessions.Default(c)
	// session.Set("status", "logined")
	// session.Save()
	c.Redirect(302, "/")
}

func getGoogleAuthHandler(c *gin.Context) {
	code := c.Query("code")
	authCode = code
	fmt.Printf("Code: %s\n", authCode)
	client := requestGoogleAuth(authCode)
	tasksService := initialGoogleTasksService(client)
	showGoogleTasksList(tasksService)
	oauthService := initialGoogleOAuthService(client)
	// showGoogleUserInfo(*oauthService)
	userInfo := getGoogleUserInfo(*oauthService)
	addLoginCookie(userInfo, c)
	// session := sessions.Default(c)
	// session.Set("status", "logined")
	// session.Save()
	c.Redirect(302, "/")
}

func getMicrosoftAuthHandler(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Printf("[NOTE] session: %v\n", session.Get("id"))
	session.Set("id", "logined")
	session.Save()
	fmt.Printf("[NOTE] session: %v\n", session.Get("id"))
	c.Redirect(302, "/")
}
