package beer

import (
	"errors"

	"github.com/landru29/api-go/model"
	uuid "github.com/nu7hatch/gouuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RecipePost define to model of a new created recipe
type RecipePost struct {
	Name string  `bson:"name"          json:"name"`
	Date float64 `bson:"date"          json:"date"`
}

// Recipe define a quizz question
type Recipe struct {
	model.Dater
	Name  string        `bson:"name"          json:"name"`
	Date  float64       `bson:"date"          json:"date"`
	ID    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	User  []string      `bson:"user" json:"-"`
	Steps []Step        `bson:"steps"         json:"steps,omitempty"`
}

// Unit defines a unit
type Unit struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
}

// Physical defines a physical
type Physical struct {
	Value float64 `bson:"value" json:"value"`
	Unit  Unit    `bson:"unit" json:"unit"`
}

// Step defines a brewing step
type Step struct {
	Name        string       `bson:"name" json:"name"`
	Lasting     float64      `bson:"lasting" json:"lasting"`
	Temperature Physical     `bson:"temperature" json:"temperature"`
	UUID        string       `bson:"_uuid" json:"_uuid"`
	Ingredients []Ingredient `bson:"ingredients" json:"ingredients,omitempty"`
}

// Ingredient defines a brewing ingredient
type Ingredient struct {
	UUID     string   `bson:"_uuid" json:"_uuid"`
	Name     string   `bson:"name" json:"name"`
	Type     string   `bson:"type" json:"type"`
	Quantity Physical `bson:"qty" json:"qty"`
	// Malt
	Yield         float64 `bson:"yield,omitempty" json:"yield,omitempty"`
	Color         float64 `bson:"color,omitempty" json:"color,omitempty"`
	RecommendMash bool    `bson:"recommendMash,omitempty" json:"recommendMash,omitempty"`
	RGB           string  `bson:"_rgb,omitempty" json:"_rgb,omitempty"`
	// Hops
	Alpha float64 `bson:"alpha,omitempty" json:"alpha,omitempty"`
	Form  string  `bson:"form,omitempty" json:"form,omitempty"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
	return db.C("beerrecipes")
}

// Save function save a recipe
func (data Recipe) Save(db *mgo.Database) (result Recipe, info *mgo.ChangeInfo, err error) {
	data.UpdateDates()
	if data.ID != "" {
		info, err = getCollection(db).UpsertId(data.ID, bson.M{"$set": data})
	} else {
		err = getCollection(db).Insert(data)
	}
	result = data
	return
}

// DeleteByID remove a recipe by its ID
func DeleteByID(db *mgo.Database, ID string, userID string) (err error) {
	err = getCollection(db).Remove(
		bson.M{
			"_id": bson.ObjectIdHex(ID),
			"user": bson.M{
				"$in": []string{
					userID,
				},
			},
		})
	return
}

// DeleteStepByID delete a step from a recipe
func DeleteStepByID(db *mgo.Database, recipeID string, stepID string, userID string) (err error) {
	if recipe, err := GetRecipe(db, recipeID, userID); err == nil {
		steps := recipe.Steps
		recipe.Steps = []Step{}
		for _, step := range steps {
			if step.UUID != stepID {
				recipe.Steps = append(recipe.Steps, step)
			}
		}
		_, _, err = recipe.Save(db)
	}
	return
}

// DeleteIngredientByID delete an ingredient from the step
func DeleteIngredientByID(db *mgo.Database, recipeID string, stepID string, ingredientID string, userID string) (err error) {
	if recipe, err := GetRecipe(db, recipeID, userID); err == nil {
		if step, ok := findStepByID(recipe.Steps, stepID); ok {
			ingredients := step.Ingredients
			step.Ingredients = []Ingredient{}
			for _, ingredient := range ingredients {
				if ingredient.UUID != ingredientID {
					step.Ingredients = append(step.Ingredients, ingredient)
				}
			}
			_, _, err = recipe.Save(db)
		} else {
			err = errors.New("Missing step")
		}
	}
	return
}

// AddStep add a new step in the recipe
func AddStep(db *mgo.Database, recipeID string, step Step, userID string) (result Step, err error) {
	if recipe, err := GetRecipe(db, recipeID, userID); err == nil {
		if UUID, err := uuid.NewV4(); err == nil {
			step.UUID = UUID.String()
			recipe.Steps = append(recipe.Steps, step)
			_, _, err = recipe.Save(db)
			result = step
		}
	}
	return
}

// AddIngredient add a new ingredient in the step
func AddIngredient(db *mgo.Database, recipeID string, stepID string, ingredient Ingredient, userID string) (result Ingredient, err error) {
	if recipe, err := GetRecipe(db, recipeID, userID); err == nil {
		if step, ok := findStepByID(recipe.Steps, stepID); ok {
			if UUID, err := uuid.NewV4(); err == nil {
				ingredient.UUID = UUID.String()
				step.Ingredients = append(step.Ingredients, ingredient)
				_, _, err = recipe.Save(db)
				result = ingredient
			}
		} else {
			err = errors.New("No step found")
		}
	}
	return
}

// GetRecipe find a unique recipe by ID
func GetRecipe(db *mgo.Database, ID string, userID string) (result Recipe, err error) {
	err = getCollection(db).Find(
		bson.M{
			"_id": bson.ObjectIdHex(ID),
			"user": bson.M{
				"$in": []string{
					userID,
				},
			},
		}).One(&result)
	return
}

// GetAllRecipes function get all Recipes
func GetAllRecipes(db *mgo.Database, skip int, count int) (results []Recipe, err error) {
	results = []Recipe{}
	err = getCollection(db).Find(bson.M{}).Skip(skip).Limit(count).All(&results)
	return
}

// GetAllRecipesByUser function get all Recipes for a given user
func GetAllRecipesByUser(db *mgo.Database, userID string, skip int, count int) (results []Recipe, err error) {
	results = []Recipe{}
	err = getCollection(db).Find(
		bson.M{
			"user": bson.M{
				"$in": []string{
					userID,
				},
			},
		}).Skip(skip).Limit(count).All(&results)
	return
}

func findStepByID(steps []Step, ID string) (step Step, ok bool) {
	ok = false
	for _, curStep := range steps {
		if curStep.UUID == ID {
			ok = true
			step = curStep
			return
		}
	}
	return
}
