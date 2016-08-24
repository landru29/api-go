package routes

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/landru29/api-go/middleware"
    "github.com/landru29/api-go/model/quizz"
    "github.com/landru29/api-go/mongo"
    "github.com/spf13/viper"
)

func apiBaseURL() string {
    port := viper.GetString("api_port")
    return viper.GetString("api_protocol") + "://" + viper.GetString("api_host") +
        map[bool]string{true: ":" + port, false: ""}[len(port) > 0] + "/"
}

func loginRedirect(c *gin.Context, id string) {
    code := c.DefaultQuery("code", "")
    fmt.Println(id)
    fmt.Println(code)

    c.String(200, "Got it!")
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
            "host": apiBaseURL(),
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
