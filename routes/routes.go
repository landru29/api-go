package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/landru29/api-go/authorization"
	"github.com/landru29/api-go/middleware"
	"github.com/landru29/api-go/mongo"
)

// GetEmail get the email address from third party
func GetEmail(c *gin.Context) (email string, ok bool) {
	email = ""
	user, ok := c.Get("user")
	if ok == true {
		email = user.(authorization.Profile).Email
	}
	return
}

// GetID get the user ID from jwt
func GetID(c *gin.Context) (ID string, ok bool) {
	ID = ""
	user, ok := c.Get("user")
	if ok == true {
		ID = user.(authorization.Profile).ID
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
