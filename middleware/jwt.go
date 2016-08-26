package middleware

import (
	"regexp"

	mgo "gopkg.in/mgo.v2"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/token"
	"github.com/spf13/viper"
)

// MyCustomClaims is the jwt model for the token
type MyCustomClaims struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
	jwt.StandardClaims
}

// JwtMiddleware get pagination information from the query
func JwtMiddleware(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtStr := c.Request.Header.Get("Authorization")
		re := regexp.MustCompile(`(?i)^jwt\s(.*)$`)
		matching := re.FindStringSubmatch(jwtStr)
		if len(matching) > 0 {
			tokenString := matching[1]
			//fmt.Println(tokenString)
			tokenJwt, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("jwt_secret")), nil
			})
			if err != nil {
				c.Next()
				return
			}
			if claims, ok := tokenJwt.Claims.(*MyCustomClaims); ok && tokenJwt.Valid {
				dbToken, err := token.FindToken(db, claims.Token)
				if err != nil {
					c.Next()
					return
				}
				c.Set("token", dbToken)
				c.Next()
				return
			}
			c.Next()
			return
		}
		c.Next()
		return
	}
}
