package routes

import (
	"fmt"

	"github.com/landru29/gin-passport-facebook"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	mgo "gopkg.in/mgo.v2"
)

func prepareFaceBook() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  apiBaseURL() + "auth/facebook/callback",
		ClientID:     viper.GetString("facebookAuth.clientId"),
		ClientSecret: viper.GetString("facebookAuth.clientSecret"),
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
}

func handleFacebook(router *gin.Engine, database *mgo.Database) {
	authFacebook := router.Group("/auth/facebook")
	opts := prepareFaceBook()
	GinPassportFacebook.Routes(opts, authFacebook)
	authFacebook.GET("/callback", GinPassportFacebook.Middleware(), func(c *gin.Context) {
		userAuth, err := GinPassportFacebook.GetProfile(c)
		if (userAuth == nil) || (err != nil) {
			c.AbortWithStatus(500)
			return
		}

		accessToken, _ := c.Get("client-token")
		fmt.Println(accessToken.(oauth2.Token))

		userDb, errDb := user.FindUser(database, userAuth.Email)
		if errDb == nil {
			userDb.Name = userAuth.Name
			userDb.Facebook.Code = c.DefaultQuery("code", "")
			userDb.Facebook.ID = userAuth.Id
			userDb.Google.Code = ""
			userDb.Save(database)
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
