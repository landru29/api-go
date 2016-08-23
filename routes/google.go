package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/user"
	"github.com/landru29/gin-passport-google"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	mgo "gopkg.in/mgo.v2"
)

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
	authGoogle := router.Group("/auth/google")
	GinPassportGoogle.Routes(prepareGoogle(), authGoogle)
	authGoogle.GET("/callback", GinPassportGoogle.Middleware(), func(c *gin.Context) {
		userAuth, err := GinPassportGoogle.GetProfile(c)
		if (userAuth == nil) || (err != nil) {
			c.AbortWithStatus(500)
			return
		}

		accessToken, _ := c.Get("client-token")
		fmt.Println(accessToken.(oauth2.Token))

		userDb, errDb := user.FindUser(database, userAuth.Email)
		if errDb == nil {
			userDb.Name = userAuth.Name
			userDb.Google.Code = c.DefaultQuery("code", "")
			userDb.Google.ID = userAuth.Id
			userDb.Facebook.Code = ""
			_, _, errSave := userDb.Save(database)
			if errSave != nil {
				c.AbortWithStatus(500)
				return
			}
		} else {
			c.AbortWithStatus(500)
			return
		}
		loginRedirect(c, userAuth.Email)
	})
}
