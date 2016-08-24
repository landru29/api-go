package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	mgo "gopkg.in/mgo.v2"
)

// Google is the name of google
const Google = "google"

func prepareGoogle() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  apiBaseURL() + "auth/google/callback",
		ClientID:     viper.GetString("googleAuth.clientId"),
		ClientSecret: viper.GetString("googleAuth.clientSecret"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

}

func handleGoogle(router *gin.Engine, database *mgo.Database) {
	googleRouter := router.Group("/auth/google")
	authGoogle := prepareGoogle()

	googleRouter.GET("/login", func(c *gin.Context) {
		url := authGoogle.AuthCodeURL("")
		redirect := "redirect=" + c.Query("redirect") + "; Path=/; HttpOnly"
		c.Header("set-cookie", redirect)
		c.Redirect(http.StatusFound, url)
	})

	googleRouter.GET("/callback", func(c *gin.Context) {

		c.Request.ParseForm()

		gCode := c.Request.Form.Get("code")

		apiToken, err := authGoogle.Exchange(c, gCode)
		if apiToken == nil {
			c.Redirect(301, "/")
			return
		} else if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		profile, err := getProfile(c, authGoogle, apiToken, Google)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		uri, err := getRedirect(c, Google, apiToken, profile)
		if err != nil {
			return
		}
		c.Redirect(301, uri)

	})
}
