package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/landru29/api-go/authorization"
    "github.com/landru29/api-go/middleware"
    "github.com/landru29/api-go/model/quizz"
    "github.com/landru29/api-go/mongo"
)

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
