package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/beer"
	"github.com/landru29/api-go/mongo"
)

func beerRoutes(router *gin.Engine) {
	database := mongo.GetMongoDatabase()

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
			count, _ := c.Get("count")
			skip, _ := c.Get("skip")
			recipes, err := beer.GetAllRecipes(database, skip.(int), count.(int))
			if err != nil {

			}
			if err != nil {
				content := gin.H{"message": "Error while reading database"}
				c.JSON(http.StatusServiceUnavailable, content)
			} else {
				c.JSON(http.StatusOK, recipes)
			}
		})
	}
}
