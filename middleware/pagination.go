package middleware

import (
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
)

// PaginationMiddleware get pagination information from the query
func PaginationMiddleware() gin.HandlerFunc {
    defaultLimit, err := strconv.Atoi(viper.GetString("default_pagination_limit"))
    if err != nil {
        defaultLimit = 1
    }

    return func(c *gin.Context) {

        // compute count
        limit := c.Query("limit")
        count := c.Query("count")
        if len(limit) > 0 {
            count = limit
        }
        outCount, err := strconv.Atoi(count)
        if err != nil {
            outCount = defaultLimit
        }
        c.Set("count", outCount)

        // compute skip
        page, err := strconv.Atoi(c.Query("page"))
        if err != nil {
            page = 1
        }
        c.Set("skip", (page-1)*outCount)

        // End of middleware
        c.Next()
    }
}
