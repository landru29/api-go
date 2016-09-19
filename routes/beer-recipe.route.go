package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func beerRoutes(router *gin.Engine) {
	//database := mongo.GetMongoDatabase()

	// @SubApi Beer [/beer]
	// @SubApi Beer resource resource [/beer]
	beerRecipeGroup := router.Group("/beer")
	{
		// @Title Get recipes
		// @Description Get all paginated recipes
		// @Accept application/json
		// @Param limit query integer false "Number of recipes to return"
		// @Param page query integer false "Page to return"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router / [get]
		beerRecipeGroup.GET("/", func(c *gin.Context) {
			/*count, _ := c.Get("count")*/
			content := gin.H{"message": "Not implemented yet"}
			c.JSON(http.StatusOK, content)
		})
	}
}
