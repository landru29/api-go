package routes

import (
	"fmt"
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
		// @Param page  query integer false "Page to return"
		// @Success 200 {object} string "Success"
		// @Failure 401 {object} string "Unauthorized"
		// @Resource /beer
		// @Router / [get]
		beerRecipeGroup.GET("/", func(c *gin.Context) {
			count, _ := c.Get("count")
			skip, _ := c.Get("skip")
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				recipes, err := beer.GetAllRecipesByUser(database, userID, skip.(int), count.(int))
				if err != nil {
					content := gin.H{"message": "Error while reading database"}
					c.JSON(http.StatusServiceUnavailable, content)
				} else {
					c.JSON(http.StatusOK, recipes)
				}
			}

		})

		// @Title Add recipe
		// @Description Add a new recipe
		// @Accept application/json
		// @Param name  body string true "Name of the recipe"
		// @Param date  body string true "Date of the recipe"
		// @Param steps body string true "Array of steps"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router / [post]
		beerRecipeGroup.POST("/", func(c *gin.Context) {
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				fmt.Println(userID)
				content := gin.H{"message": "Not implemented yet"}
				c.JSON(http.StatusServiceUnavailable, content)
			}

		})

		// @Title Add recipe
		// @Description Update a recipe
		// @Accept application/json
		// @param id    path string true "Identifier of the recipe"
		// @Param name  body string true "Name of the recipe"
		// @Param date  body string true "Date of the recipe"
		// @Param steps body string true "Array of steps"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router /:id [put]
		beerRecipeGroup.PUT("/:id", func(c *gin.Context) {
			ID := c.Param("id")
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				fmt.Println(ID, userID)
				content := gin.H{"message": "Not implemented yet"}
				c.JSON(http.StatusServiceUnavailable, content)
			}

		})
	}
}
