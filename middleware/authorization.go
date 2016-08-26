package middleware

import (
	"regexp"

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
				c.Next()
				return
			}
			var oauthConf *oauth2.Config
			switch dbToken.Origin {
			case authorization.Facebook:
				oauthConf = authorization.GetFacebookOAuth()
				c.Set("oauth", *oauthConf)
				c.Set("dbToken", dbToken)
			case authorization.Google:
				oauthConf = authorization.GetGoogleOAuth()
				c.Set("oauth", *oauthConf)
				c.Set("dbToken", dbToken)
			}
			c.Next()
			return
		}
		c.Next()
		return
	}
}
