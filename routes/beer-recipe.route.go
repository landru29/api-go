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
		// @Param recipe body string true "recipe"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router / [post]
		beerRecipeGroup.POST("/", func(c *gin.Context) {
			recipe := beer.Model{}
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				if err := c.BindJSON(&recipe); err == nil {
					recipe.User = []string{userID}
					fmt.Println(recipe)
					if result, _, err := recipe.Save(database); err == nil {
						c.JSON(http.StatusOK, result)
					} else {
						content := gin.H{"message": "Could not save"}
						c.JSON(http.StatusServiceUnavailable, content)
					}
				}
			}
		})

		// @Title Add recipe
		// @Description Update a recipe
		// @Accept application/json
		// @param recipeId path string true "Identifier of the recipe"
		// @Param name     body string true "Name of the recipe"
		// @Param date     body string true "Date of the recipe"
		// @Param steps    body string true "Array of steps"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router /:recipeId [put]
		beerRecipeGroup.PUT("/:recipeId", func(c *gin.Context) {
			recipeID := c.Param("recipeId")
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				fmt.Println(recipeID, userID)
				content := gin.H{"message": "Not implemented yet"}
				c.JSON(http.StatusServiceUnavailable, content)
			}

		})

		// @Title Add step
		// @Description Add a new step to the recipe
		// @Accept application/json
		// @param recipeId path string true "Identifier of the recipe"
		// @Param step     body string true "Step"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router /:recipeId/step [post]
		beerRecipeGroup.POST("/:recipeId/step", func(c *gin.Context) {
			recipeID := c.Param("recipeId")
			//recipe := beer.Model{}
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				fmt.Println(recipeID, userID)
				content := gin.H{"message": "Not implemented yet"}
				c.JSON(http.StatusServiceUnavailable, content)
			}
		})

		// @Title Delete step
		// @Description Add a new step to the recipe
		// @Accept application/json
		// @param recipeId path string true "Identifier of the recipe"
		// @param stepId path string true "Identifier of the step"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router /:recipeId/step [post]
		beerRecipeGroup.DELETE("/:recipeId/step/:stepId", func(c *gin.Context) {
			recipeID := c.Param("recipeId")
			stepID := c.Param("stepId")
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				fmt.Println(recipeID, stepID, userID)
				content := gin.H{"message": "Not implemented yet"}
				c.JSON(http.StatusServiceUnavailable, content)
			}
		})
	}
}
