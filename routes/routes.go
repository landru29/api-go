package routes

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	"github.com/durango/gin-passport-facebook"
	"github.com/durango/gin-passport-google"
	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/middleware"
	"github.com/landru29/api-go/model/quizz"
	"github.com/landru29/api-go/model/user"
	"github.com/landru29/api-go/mongo"
	"github.com/spf13/viper"
)

func apiBaseURL() string {
	port := viper.GetString("api_port")
	return viper.GetString("api_protocol") + "://" + viper.GetString("api_host") +
		map[bool]string{true: ":" + port, false: ""}[len(port) > 0] + "/"
}

func prepareFaceBook() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  apiBaseURL() + "auth/facebook/callback",
		ClientID:     viper.GetString("facebookAuth.clientId"),
		ClientSecret: viper.GetString("facebookAuth.clientSecret"),
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
}

func prepareGoogle() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  apiBaseURL() + "auth/google/callback",
		ClientID:     viper.GetString("googleAuth.clientId"),
		ClientSecret: viper.GetString("googleAuth.clientSecret"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

}

func loginRedirect(c *gin.Context, id string) {
	code := c.DefaultQuery("code", "")
	fmt.Println(id)
	fmt.Println(code)

	c.String(200, "Got it!")
}

func handleFacebook(router *gin.Engine, database *mgo.Database) {
	authFacebook := router.Group("/auth/facebook")
	GinPassportFacebook.Routes(prepareFaceBook(), authFacebook)
	authFacebook.GET("/callback", GinPassportFacebook.Middleware(), func(c *gin.Context) {
		userAuth, err := GinPassportFacebook.GetProfile(c)
		if (userAuth == nil) || (err != nil) {
			c.AbortWithStatus(500)
			return
		}
		userDb, errDb := user.FindUser(database, userAuth.Email)
		if errDb != nil {
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
		}
		loginRedirect(c, userAuth.Email)
	})
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
		userDb, errDb := user.FindUser(database, userAuth.Email)
		if errDb != nil {
			userDb.Name = userAuth.Name
			userDb.Google.Code = c.DefaultQuery("code", "")
			userDb.Google.ID = userAuth.Id
			userDb.Facebook.Code = ""
			_, _, errSave := userDb.Save(database)
			if errSave != nil {
				c.AbortWithStatus(500)
				return
			}
		}
		loginRedirect(c, userAuth.Email)
	})
}

// DefineRoutes defines all routes
func DefineRoutes() *gin.Engine {
	database := mongo.GetMongoDatabase()
	router := gin.Default()

	// facebook
	handleFacebook(router, database)

	// google
	handleGoogle(router, database)

	// Middlewares
	router.Use(middleware.PaginationMiddleware())

	router.LoadHTMLGlob("./templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	quizzGroup := router.Group("/quizz")
	{
		quizzGroup.GET("/", func(c *gin.Context) {
			count, _ := c.Get("count")
			results, err := quizz.RandomPublished(database, count.(int), 10)
			if err != nil {
				content := gin.H{"message": "Error while reading database"}
				c.JSON(http.StatusServiceUnavailable, content)
			} else {
				c.JSON(http.StatusOK, results)
			}
		})
	}

	return router
}
