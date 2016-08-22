package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GoogleModel define a google connect
type GoogleModel struct {
	ID   string `bson:"id"`
	Code string `bson:"token"`
}

// FacebookModel define a facebook connect
type FacebookModel struct {
	ID   string `bson:"id"`
	Code string `bson:"token"`
}

// Model define a user
type Model struct {
	ID           bson.ObjectId `bson:"_id"`
	Email        string        `bson:"email"`
	Password     string        `bson:"password"`
	Applications []string      `bson:"applications"`
	Role         string        `bson:"role"`
	Active       bool          `bson:"active"`
	Verified     bool          `bson:"verified"`
	Google       GoogleModel   `bson:"google"`
	Facebook     FacebookModel `bson:"facebook"`
	CreatedAt    string        `bson:"createdAt"`
	UpdatedAt    string        `bson:"updatedAt"`
	Name         string        `bson:"name"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
	return db.C("users")
}

// FindUser find a user
func FindUser(db *mgo.Database, email string) (result Model, err error) {
	err = getCollection(db).Find(bson.M{"email": email}).One(&result)
	return
}

// Save save a user
func (data Model) Save(db *mgo.Database) (result Model, info *mgo.ChangeInfo, err error) {
	info, err = getCollection(db).Upsert(bson.M{"email": data.Email}, bson.M{"$set": data})
	result = data
	return
}
