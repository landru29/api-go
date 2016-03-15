package mongo

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"github.com/spf13/viper"
	"github.com/landru29/api-go/helpers/config"
)

const (
	databaseName          = "api"
	collectionQuizz       = "quizzs"
	collectionApplication = "applications"
	collectionBeerRecipe  = "beerrecipes"
	collectionUser        = "users"
)

type MongoStore struct {
	clQuizz       *mgo.Collection
	clApplication *mgo.Collection
	clBeerRecipe  *mgo.Collection
	clUser        *mgo.Collection
}

var _currentSession *mgo.Session
var  _instance *MongoStore

func GetSession() (*mgo.Session) {
	return _currentSession
}

func GetInstance() (*MongoStore) {
	return _instance;
}

func Connect() (session *mgo.Session, err error) {
	config.Load()
	host := getDbParam("host")
	port := getDbParam("port")
	name := getDbParam("name")
	user := getDbParam("user")
	password := getDbParam("password")

	mongooseConnectionChain := "mongodb://" +
	map[bool]string{true: user, false:""} [len(user)>0] +
	map[bool]string{true: ":" + password, false:""} [len(password)>0] +
	map[bool]string{true: "@", false:""} [len(user)>0] +
	host +
	map[bool]string{true: ":" + port, false:""} [len(port)>0]

	session, err =  mgo.Dial(mongooseConnectionChain)
	_currentSession = session

	if err != nil {
		panic(err)
	}
	defer session.Close()

	_instance = &MongoStore{
		clQuizz:       _currentSession.DB(name).C(collectionQuizz),
		clApplication: _currentSession.DB(name).C(collectionApplication),
		clBeerRecipe:  _currentSession.DB(name).C(collectionBeerRecipe),
		clUser:        _currentSession.DB(name).C(collectionUser),
	}

	return
}

func getDbParam(key string) string {
	return viper.GetString("application.database." + key)
}