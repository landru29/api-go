package middleware

import (
    "fmt"
    "regexp"
    "time"

    "golang.org/x/oauth2"

    mgo "gopkg.in/mgo.v2"

    "github.com/gin-gonic/gin"
    "github.com/landru29/api-go/authorization"
)

// AuthorizationMiddleware get pagination information from the query
func AuthorizationMiddleware(db *mgo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        jwtStr := c.Request.Header.Get("Authorization")
        re := regexp.MustCompile(`(?i)^jwt\s(.*)$`)
        matching := re.FindStringSubmatch(jwtStr)
        if len(matching) > 0 {
            tokenString := matching[1]
            dbToken, err := authorization.DecodeToken(tokenString, db)
            if err != nil {
                fmt.Println("Bad JWT")
                c.Next()
                return
            }
            myToken := oauth2.Token{
                AccessToken:  dbToken.AccessToken,
                TokenType:    dbToken.TokenType,
                RefreshToken: dbToken.RefreshToken,
                Expiry:       time.Unix(dbToken.Expiry, 0),
            }
            var oauthConf *oauth2.Config
            switch dbToken.Origin {
            case authorization.Facebook:
                oauthConf = authorization.GetFacebookOAuth()
                c.Set("oauthConfig", *oauthConf)
                c.Set("oauthToken", myToken)
                email, err := authorization.GetEmail(c, oauthConf, &myToken)
                c.Set("verified", len(email) > 0)
                if err == nil {
                    c.Set("email", email)
                }
            case authorization.Google:
                oauthConf = authorization.GetGoogleOAuth()
                c.Set("oauthConfig", *oauthConf)
                c.Set("oauthToken", myToken)
                email, err := authorization.GetEmail(c, oauthConf, &myToken)
                c.Set("verified", len(email) > 0)
                if err == nil {
                    c.Set("email", email)
                }
            }
            c.Next()
            return
        }
        c.Next()
        return
    }
}
