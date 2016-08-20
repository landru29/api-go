package routes

import "github.com/gin-gonic/gin"

// DefineRoutes defines all routes
func DefineRoutes() *gin.Engine {
	router := gin.Default()

	quizz := router.Group("/quizz")
	{
		quizz.GET("/", func(c *gin.Context) {
			content := gin.H{"Hello": "World"}
			c.JSON(200, content)
		})
	}

	return router
}
