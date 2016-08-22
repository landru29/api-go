package routes

import (
    "net/http"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/facebook"
    "golang.org/x/oauth2/google"

    "github.com/durango/gin-passport-facebook"
    "github.com/durango/gin-passport-google"
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

// DefineRoutes defines all routes
func DefineRoutes() *gin.Engine {
    database := mongo.GetMongoDatabase()
    router := gin.Default()

    // facebook
    authFacebook := router.Group("/auth/facebook")
    GinPassportFacebook.Routes(prepareFaceBook(), authFacebook)
    authFacebook.GET("/callback", GinPassportFacebook.Middleware(), func(c *gin.Context) {
        user, err := GinPassportFacebook.GetProfile(c)
        if user == nil || err != nil {
            c.AbortWithStatus(500)
            return
        }

        c.String(200, "Got it, Facebook!")
    })

    // google
    authGoogle := router.Group("/auth/google")
    GinPassportGoogle.Routes(prepareGoogle(), authGoogle)
    authGoogle.GET("/callback", GinPassportGoogle.Middleware(), func(c *gin.Context) {
        user, err := GinPassportGoogle.GetProfile(c)
        if user == nil || err != nil {
            c.AbortWithStatus(500)
            return
        }

        c.String(200, "Got it, Google!")
    })

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
