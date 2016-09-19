package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/quizz"
	"github.com/landru29/api-go/mongo"
)

func quizzRoutes(router *gin.Engine) {
	database := mongo.GetMongoDatabase()

	// @SubApi Quizz [/quizz]
	// @SubApi Roller Derby Quizz resource [/quizz]
	quizzGroup := router.Group("/quizz")
	{
		// @Title Read
		// @Description Get Random questions
		// @Accept application/json
		// @Param limit query integer false "Max questions"
		// @Success 200 {object} string "Success"
		// @Resource /quizz
		// @Router / [get]
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
}
