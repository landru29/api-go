package quizz

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// Model define a quizz question
type Model struct {
    ID    bson.ObjectId `bson:"_id,omitempty"`
    Email string        `bson:"email"`
    Name  string        `bson:"name"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
    return db.C("me")
}

// Save function save a question
func (data Model) Save(db *mgo.Database) (result Model, info *mgo.ChangeInfo, err error) {
    _, err = Find(db, data.Email)
    if err == mgo.ErrNotFound {
        err = getCollection(db).Insert(data)
        result = data
        return
    }
    if err != nil {
        info, err = getCollection(db).Upsert(bson.M{"email": data.Email}, bson.M{"$set": data})
    }
    result = data
    return
}

// Find function find an element
func Find(db *mgo.Database, email string) (result Model, err error) {
    result = Model{}
    err = getCollection(db).Find(bson.M{"email": email}).One(&result)
    return
}
