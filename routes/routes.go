package routes

import (
    "errors"
    "net/http"

    "golang.org/x/oauth2"

    "github.com/gin-gonic/gin"
    "github.com/landru29/api-go/authorization"
    "github.com/landru29/api-go/middleware"
    "github.com/landru29/api-go/model/quizz"
    "github.com/landru29/api-go/mongo"
    "github.com/itsjamie/gin-cors"
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
        Origins:        "*",
        Methods:        "GET, PUT, POST, DELETE",
        RequestHeaders: "Origin, Authorization, Content-Type",
        ExposedHeaders: "",
        Credentials: true,
        ValidateHeaders: false,
    }))

    router.LoadHTMLGlob("./templates/*")

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "redirect": authorization.APIBaseURL() + "doc",
        })
    })

    router.Static("/doc", "./swagger")


    // @SubApi Quizz [/quizz]
    // @SubApi Roller Derby Quizz resource [/quizz]
    quizzGroup := router.Group("/quizz")
    {
        // @Title Read
        // @Description Get Random questions
        // @Accept application/json
        // @Param limit query int false "Max questions"
        // @Success 200 {object} string "Success"
        // @Resource /quizz
        // @Router / [get]
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

    meGroup := router.Group("/me")
    {
        meGroup.POST("/", func(c *gin.Context) {
            user, ok := c.Get("user")
            if ok == true {
                c.JSON(http.StatusOK, gin.H{
                    "message": "OK",
                    "email":   user.(authorization.Profile).Email,
                })
            } else {
                c.JSON(http.StatusOK, gin.H{
                    "message": "ko",
                })
            }
        })
    }

    return router
}
