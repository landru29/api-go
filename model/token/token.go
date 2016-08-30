package token

import (
    mgo "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// Model define a user
// swagger:model Token
type Model struct {
    ID           bson.ObjectId `bson:"_id,omitempty"`
    TokenType    string        `bson:"tokenType"`
    AccessToken  string        `bson:"accressToken"`
    RefreshToken string        `bson:"refreshToken"`
    Origin       string        `bson:"origin"`
    Email        string        `bson:"email"`
    Expiry       int64         `bson:"expiry"`
    FirstName    string        `bson:"firstName"`
    LastName     string        `bson:"lastName"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
    return db.C("token")
}

// FindToken find a user
func FindToken(db *mgo.Database, ID string) (result Model, err error) {
    err = getCollection(db).FindId(bson.ObjectIdHex(ID)).One(&result)
    return
}

// Save save a user
func (data Model) Save(db *mgo.Database) (result Model, err error) {
    result = Model{}
    tmp := Model{}
    err = getCollection(db).Find(bson.M{"email": data.Email}).One(&tmp)
    if err != nil {
        if err == mgo.ErrNotFound {
            data.ID = bson.NewObjectId()
            err = getCollection(db).Insert(data)
            result = data
        }
        return
    }
    data.ID = tmp.ID
    result = data

    err = getCollection(db).Update(bson.M{"email": data.Email}, data)

    return
}
