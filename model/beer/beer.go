package beer

import (
	"github.com/landru29/api-go/model"
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

// StepPost define to model of a new created step
type StepPost struct {
	Name        string   `bson:"name" json:"name"`
	Lasting     float64  `bson:"lasting" json:"lasting"`
	Temperature Physical `bson:"temperature" json:"temperature"`
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
