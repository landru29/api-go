package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/authorization"
)

func meRoutes(router *gin.Engine) {
	// @SubApi Me [/me]
	// @SubApi Information about current User [/me]
	meGroup := router.Group("/me")
	{
		// @Title Read
		// @Description Get informatoin about current user
		// @Accept application/json
		// @Success 200 {object} string "Success"
		// @Failure 404 {object} string "Error"
		// @Resource /me
		// @Router / [get]
		meGroup.GET("/", func(c *gin.Context) {
			user, ok := c.Get("user")
			if ok == true {
				c.JSON(http.StatusOK, gin.H{
					"email": user.(authorization.Profile).Email,
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "ko",
				})
			}
		})
	}
}
