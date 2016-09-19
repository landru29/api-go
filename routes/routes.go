package routes

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/landru29/api-go/authorization"
	"github.com/landru29/api-go/middleware"
	"github.com/landru29/api-go/mongo"
)

// GetEmail get the email address from third party
func GetEmail(c *gin.Context) (email string, err error) {
	err = errors.New("No OAuth token")
	email = ""
	oAuthConfig, ok := c.Get("oauthConfig")
	if ok == true {
		oAuthToken, ok := c.Get("oauthToken")
		if ok == true {
			authConfig := oAuthConfig.(oauth2.Config)
			authToken := oAuthToken.(oauth2.Token)
			email, err = authorization.GetEmail(c, &authConfig, &authToken)
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
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.LoadHTMLGlob("./templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"redirect": authorization.APIBaseURL() + "doc",
		})
	})

	router.Static("/doc", "./swagger")

	quizzRoutes(router)
	meRoutes(router)
	beerRoutes(router)

	return router
}
