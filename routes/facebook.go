package routes

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

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
    facebookRouter := router.Group("/auth/facebook")
    authFacebook := prepareFaceBook()

    facebookRouter.GET("/login", func(c *gin.Context) {
        url := authFacebook.AuthCodeURL("")
        c.Redirect(http.StatusFound, url)
    })

    facebookRouter.GET("/callback", func(c *gin.Context) {
        c.Request.ParseForm()

        fbCode := c.Request.Form.Get("code")

        token, err := authFacebook.Exchange(c, fbCode)

        if token == nil {
            c.Redirect(301, "/")
            return
        } else if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        client := authFacebook.Client(c, token)

        resp, err := client.Get("https://graph.facebook.com/v2.2/me?fields=id,name,email,picture,first_name,last_name")
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        defer resp.Body.Close()
        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        var userInformation Profile
        err = json.Unmarshal(contents, &userInformation)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        userDb, errDb := user.FindUser(database, userInformation.Email)
        if errDb == nil {
            userDb.Name = userInformation.Name
            userDb.Facebook.Code = fbCode
            userDb.Facebook.ID = userInformation.ID
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
        loginRedirect(c, userInformation.Email)
    })
}
