package routes

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/authorization"
	"github.com/landru29/api-go/middleware"
	"github.com/landru29/api-go/model/quizz"
	"github.com/landru29/api-go/model/token"
	"github.com/landru29/api-go/mongo"
)

// GetEmail get the email address from third party
func GetEmail(c *gin.Context) (email string, err error) {
	err = errors.New("No OAuth token")
	email = ""
	oAuth, ok := c.Get("oauth")
	if ok == true {
		dbToken, ok := c.Get("dbToken")
		if ok == true {
			auth := oAuth.(oauth2.Config)
			oauthToken := oauth2.Token{
				TokenType:    dbToken.(token.Model).TokenType,
				AccessToken:  dbToken.(token.Model).AccessToken,
				RefreshToken: dbToken.(token.Model).RefreshToken,
				Expiry:       time.Unix(dbToken.(token.Model).Expiry, 0),
			}
			email, err = authorization.GetEmail(c, &auth, &oauthToken)
		}
	}
	return
}

// DefineRoutes defines all routes
func DefineRoutes() *gin.Engine {
	database := mongo.GetMongoDatabase()
	router := gin.Default()

	// facebook
	authorization.HandleFacebook(router, database)

	// google
	authorization.HandleGoogle(router, database)

	// Middlewares
	router.Use(middleware.PaginationMiddleware())
	router.Use(middleware.AuthorizationMiddleware(database))

	router.LoadHTMLGlob("./templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"host": authorization.APIBaseURL(),
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
