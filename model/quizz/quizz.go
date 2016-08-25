package quizz

import (
    "math/rand"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// Choice define a choice in the response
type Choice struct {
    Text    string `bson:"text"`
    Scoring int    `bson:"scoring"`
}

// Model define a quizz question
type Model struct {
    ID           bson.ObjectId `bson:"_id,omitempty"`
    Explaination string        `bson:"explaination"`
    Image        string        `bson:"image"`
    Level        int           `bson:"level"`
    Published    bool          `bson:"published"`
    Tags         string        `bson:"tags"`
    Text         string        `bson:"text"`
    Choices      []Choice      `bson:"choices"`
    CreatedAt    string        `bson:"createdAt"`
    UpdatedAt    string        `bson:"updatedAt"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
    return db.C("quizzs")
}

// Save function save a question
func (data Model) Save(db *mgo.Database) (result Model, info *mgo.ChangeInfo, err error) {
    if data.ID != "" {
        info, err = getCollection(db).UpsertId(data.ID, bson.M{"$set": data})
    } else {
        err = getCollection(db).Insert(data)
    }
    result = data
    return
}

// GetAllPublished function get all published questions
func GetAllPublished(db *mgo.Database) (results []Model, err error) {
    results = []Model{}
    err = getCollection(db).Find(bson.M{"published": true}).All(&results)
    return
}

// RandomPublished function get randomize questions in db
func RandomPublished(db *mgo.Database, count int, level int) (results []Model, err error) {
    results = []Model{}
    // get db count
    total, err := CountPublished(db, level)
    if err != nil {
        return
    }
    // build the list
    var buffer []int
    for i := 0; i < total; i++ {
        buffer = append(buffer, i)
    }
    // random Published
    for i := 0; i < count; i++ {
        if len(buffer) > 0 {
            var element Model
            index := rand.Intn(len(buffer))
            skip := buffer[index]
            buffer = append(buffer[:index], buffer[index+1:]...)
            getCollection(db).Find(bson.M{"published": true, "level": level}).Skip(skip).Limit(1).One(&element)
            results = append(results, element)
        }
    }
    return
}

// Find function find an element
func Find(db *mgo.Database, id string) (result Model, err error) {
    result = Model{}
    err = getCollection(db).FindId(id).One(&result)
    return
}

// CountAll count all the questions
func CountAll(db *mgo.Database) (result int, err error) {
    result, err = getCollection(db).Count()
    return
}

// CountPublished count all published questions
func CountPublished(db *mgo.Database, level int) (result int, err error) {
    if level < 0 {
        result, err = getCollection(db).Find(bson.M{"published": true}).Count()
    } else {
        result, err = getCollection(db).Find(bson.M{"published": true, "level": level}).Count()
    }
    return
}
