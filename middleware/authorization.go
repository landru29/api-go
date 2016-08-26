package middleware

import (
    "regexp"

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
            c.Set("token", dbToken)
            c.Next()
            return
        }
        c.Next()
        return
    }
}
