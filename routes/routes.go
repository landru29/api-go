package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/middleware"
	"github.com/landru29/api-go/model/quizz"
	"github.com/landru29/api-go/mongo"
)

// DefineRoutes defines all routes
func DefineRoutes() *gin.Engine {
	database := mongo.GetMongoDatabase()
	router := gin.Default()
	router.Use(middleware.PaginationMiddleware())

	quizzGroup := router.Group("/quizz")
	{
		quizzGroup.GET("/", func(c *gin.Context) {
			count, _ := c.Get("count")
			results, err := quizz.RandomPublished(database, count.(int), 10)
			if err != nil {
				content := gin.H{"message": "Error while reading database"}
				c.JSON(503, content)
			} else {
				c.JSON(200, results)
			}
		})
	}

	return router
}
