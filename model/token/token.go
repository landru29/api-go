package token

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Model define a user
type Model struct {
	ID           bson.ObjectId `bson:"_id"`
	TokenType    string        `bson:"tokenType"`
	AccessToken  string        `bson:"accressToken"`
	RefreshToken string        `bson:"refreshToken"`
	Origin       string        `bson:"origin"`
	Expiry       int64         `bson:"expiry"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
	return db.C("token")
}

// FindToken find a user
func FindToken(db *mgo.Database, ID string) (result Model, err error) {
	err = getCollection(db).FindId(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&result)
	return
}

// Save save a user
func (data Model) Save(db *mgo.Database) (result Model, err error) {
	/*if data.ID == nil {
		data.ID = bson.NewObjectId().Hex()
		_, err = getCollection(db).Insert(data)
	} else {
		err := getCollection(db).UpdateId(data.ID, bson.M{"$set": data})
	}*/
	data.ID = bson.NewObjectId()
	err = getCollection(db).Insert(data)
	result = data
	return
}
