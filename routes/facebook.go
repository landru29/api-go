package routes

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	mgo "gopkg.in/mgo.v2"
)

// Facebook is the name of facebook
const Facebook = "facebook"

func prepareFacebook() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  apiBaseURL() + "auth/facebook/callback",
		ClientID:     viper.GetString("facebookAuth.clientId"),
		ClientSecret: viper.GetString("facebookAuth.clientSecret"),
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
}

func handleFacebook(router *gin.Engine, database *mgo.Database) {
	facebookRouter := router.Group("/auth/facebook")
	authFacebook := prepareFacebook()

	facebookRouter.GET("/login", func(c *gin.Context) {
		url := authFacebook.AuthCodeURL("")
		redirect := "redirect=" + c.Query("redirect") + "; Path=/; HttpOnly"
		c.Header("set-cookie", redirect)
		c.Redirect(http.StatusFound, url)
	})

	facebookRouter.GET("/callback", func(c *gin.Context) {

		c.Request.ParseForm()

		fbCode := c.Request.Form.Get("code")

		apiToken, err := authFacebook.Exchange(c, fbCode)
		if apiToken == nil {
			c.Redirect(301, "/")
			return
		} else if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		email, err := getEmail(c, authFacebook, apiToken, Facebook)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		uri, err := getRedirect(c, Facebook, apiToken, email)
		if err != nil {
			return
		}
		c.Redirect(301, uri)
	})
}
