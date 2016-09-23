package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/beer"
	"github.com/landru29/api-go/mongo"
	uuid "github.com/nu7hatch/gouuid"
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
					c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Error while reading database"})
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
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				recipePost := beer.RecipePost{}
				if err := c.BindJSON(&recipePost); err == nil {
					recipe := beer.Recipe{
						User: []string{userID},
						Name: recipePost.Name,
						Date: recipePost.Date,
					}
					if result, _, err := recipe.Save(database); err == nil {
						c.JSON(http.StatusOK, result)
					} else {
						c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save"})
					}
				}
			}
		})

		// @Title Update recipe
		// @Description Update a recipe
		// @Accept application/json
		// @param recipeId path string true "Identifier of the recipe"
		// @Param recipe   body string true "recipe"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router /:recipeId [put]
		// @Deprecated
		beerRecipeGroup.PUT("/:recipeId", func(c *gin.Context) {
			recipeID := c.Param("recipeId")
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				fmt.Println(recipeID, userID)
				c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Update recipe not implemented yet"})
			}

		})

		// @Title Delete recipe
		// @Description Delete a recipe
		// @Accept application/json
		// @param recipeId path string true "Identifier of the recipe"
		// @Success 200 {object} string "Success"
		// @Resource /beer
		// @Router /:recipeId [delete]
		beerRecipeGroup.DELETE("/:recipeId", func(c *gin.Context) {
			recipeID := c.Param("recipeId")
			if userID, ok := GetID(c); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
			} else {
				if err := beer.DeleteByID(database, recipeID, userID); err == nil {
					c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted"})
				} else {
					c.JSON(http.StatusNotFound, gin.H{"message": "Could not delete the recipe"})
				}
			}
		})

		// @SubApi [beer-step]
		// @Resource /beer
		// @Router /:recipeId/step
		beerRecipeStepGroup := beerRecipeGroup.Group("/:recipeId/step")
		{
			// @Title Add step
			// @Description Add a new step to the recipe
			// @Accept application/json
			// @param recipeId path string true "Identifier of the recipe"
			// @Param step     body string true "Step"
			// @Success 200 {object} string "Success"
			// @Resource beer-step
			// @Router / [post]
			beerRecipeStepGroup.POST("/", func(c *gin.Context) {
				recipeID := c.Param("recipeId")
				if userID, ok := GetID(c); !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
				} else {
					if recipe, err := beer.GetRecipe(database, recipeID, userID); err == nil {
						stepPost := beer.StepPost{}
						if err := c.BindJSON(&stepPost); err == nil {
							if UUID, err := uuid.NewV4(); err == nil {
								step := beer.Step{
									Name:        stepPost.Name,
									Lasting:     stepPost.Lasting,
									Temperature: stepPost.Temperature,
									UUID:        UUID.String(),
								}
								recipe.Steps = append(recipe.Steps, step)
								if result, _, err := recipe.Save(database); err == nil {
									c.JSON(http.StatusOK, result)
								} else {
									c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save"})
								}
							} else {
								c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong. Please try again."})
							}
						}
					} else {
						c.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
					}
				}
			})

			// @Title Update step
			// @Description Update a step in the recipe
			// @Accept application/json
			// @param recipeId path string true "Identifier of the recipe"
			// @param stepId path string true "Identifier of the step"
			// @Param step     body string true "Step"
			// @Success 200 {object} string "Success"
			// @Resource beer-step
			// @Router /:stepId [put]
			// @Deprecated
			beerRecipeStepGroup.PUT("/:stepId", func(c *gin.Context) {
				recipeID := c.Param("recipeId")
				stepID := c.Param("stepId")
				if userID, ok := GetID(c); !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
				} else {
					fmt.Println(recipeID, stepID, userID)
					content := gin.H{"message": "Update step not implemented yet"}
					c.JSON(http.StatusServiceUnavailable, content)
				}
			})

			// @Title Delete step
			// @Description Delete a step from a recipe
			// @Accept application/json
			// @param recipeId path string true "Identifier of the recipe"
			// @param stepId   path string true "Identifier of the step"
			// @Success 200 {object} string "Success"
			// @Resource beer-step
			// @Router /:stepId [delete]
			beerRecipeStepGroup.DELETE("/:stepId", func(c *gin.Context) {
				recipeID := c.Param("recipeId")
				stepID := c.Param("stepId")
				if userID, ok := GetID(c); !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
				} else {
					if err := beer.DeleteStepByID(database, recipeID, stepID, userID); err == nil {
						c.JSON(http.StatusOK, gin.H{"message": "Step deleted"})
					} else {
						c.JSON(http.StatusNotFound, gin.H{"message": "Could not delete the step"})
					}
				}

			})

			// @SubApi [beer-step-ingredient]
			// @Resource beer-step
			// @Router /:stepId/ingredient
			beerRecipeStepIngredientGroup := beerRecipeStepGroup.Group("/:stepId/ingredient")
			{

				// @Title Add Ingredient
				// @Description Add a new ingredient in the step
				// @Accept application/json
				// @param recipeId path string true "Identifier of the recipe"
				// @param stepId   path string true "Identifier of the step"
				// @Param ingredient     body string true "Ingredient"
				// @Success 200 {object} string "Success"
				// @Resource beer-step-ingredient
				// @Router / [post]
				// @Deprecated
				beerRecipeStepIngredientGroup.POST("/", func(c *gin.Context) {
					recipeID := c.Param("recipeId")
					stepID := c.Param("stepId")
					if userID, ok := GetID(c); !ok {
						c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
					} else {
						fmt.Println(recipeID, stepID, userID)
						content := gin.H{"message": "Add ingredient not implemented yet"}
						c.JSON(http.StatusServiceUnavailable, content)
					}
				})

				// @Title Update ingredient
				// @Description Update an ingredient in the step
				// @Accept application/json
				// @param recipeId     path string true "Identifier of the recipe"
				// @param stepId       path string true "Identifier of the step"
				// @param ingredientId path string true "Identifier of the ingredient"
				// @Param ingredient     body string true "Step"
				// @Success 200 {object} string "Success"
				// @Resource beer-step-ingredient
				// @Router /:ingredientId [put]
				// @Deprecated
				beerRecipeStepIngredientGroup.PUT("/:ingredientId", func(c *gin.Context) {
					recipeID := c.Param("recipeId")
					stepID := c.Param("stepId")
					ingredientID := c.Param("ingredientId")
					if userID, ok := GetID(c); !ok {
						c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
					} else {
						fmt.Println(recipeID, stepID, ingredientID, userID)
						content := gin.H{"message": "Update ingredient not implemented yet"}
						c.JSON(http.StatusServiceUnavailable, content)
					}
				})

				// @Title Delete step
				// @Description Delete an ingredient from a step
				// @Accept application/json
				// @param recipeId     path string true "Identifier of the recipe"
				// @param stepId       path string true "Identifier of the step"
				// @param ingredientId path string true "Identifier of the ingredient"
				// @Success 200 {object} string "Success"
				// @Resource beer-step-ingredient
				// @Router /:ingredientId [delete]
				beerRecipeStepIngredientGroup.DELETE("/:ingredientId", func(c *gin.Context) {
					recipeID := c.Param("recipeId")
					stepID := c.Param("stepId")
					ingredientID := c.Param("ingredientId")
					if userID, ok := GetID(c); !ok {
						c.JSON(http.StatusUnauthorized, gin.H{"message": "You must login before"})
					} else {
						if err := beer.DeleteIngredientByID(database, recipeID, stepID, ingredientID, userID); err == nil {
							c.JSON(http.StatusOK, gin.H{"message": "Ingredient deleted"})
						} else {
							c.JSON(http.StatusNotFound, gin.H{"message": "Could not delete the ingredient"})
						}
					}

				})

			}

		}

	}
}
