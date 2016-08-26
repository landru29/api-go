package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/mgo.v2"
)

// Google is the name of google
const Google = "google"

// GetGoogleOAuth build an OAuth object for Facebook
func GetGoogleOAuth() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  APIBaseURL() + "auth/google/callback",
		ClientID:     viper.GetString("googleAuth.clientId"),
		ClientSecret: viper.GetString("googleAuth.clientSecret"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

}

// HandleGoogle is the facebook handler for login
func HandleGoogle(router *gin.Engine, database *mgo.Database) {
	googleRouter := router.Group("/auth/google")
	authGoogle := GetGoogleOAuth()

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

		email, err := GetEmail(c, authGoogle, apiToken)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		uri, err := buildRedirect(c, Google, apiToken, email, database)
		if err != nil {
			return
		}
		c.Redirect(301, uri)

	})
}
