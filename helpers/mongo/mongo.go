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
	Quizz       *mgo.Collection
	Application *mgo.Collection
	BeerRecipe  *mgo.Collection
	User        *mgo.Collection
}

var _currentSession *mgo.Session
var  _instance MongoStore

func GetSession() (*mgo.Session) {
	return _currentSession
}

func GetInstance() (MongoStore) {
	return _instance;
}

func Connect(host, port, user, password, name string) (session *mgo.Session, err error) {
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
	//defer session.Close()

	_instance = MongoStore{
		Quizz:       _currentSession.DB(name).C(collectionQuizz),
		Application: _currentSession.DB(name).C(collectionApplication),
		BeerRecipe:  _currentSession.DB(name).C(collectionBeerRecipe),
		User:        _currentSession.DB(name).C(collectionUser),
	}

	return
}

func AutoConnect() (session *mgo.Session, err error) {
	config.Load()
	host := getDbParam("host")
	port := getDbParam("port")
	name := getDbParam("name")
	user := getDbParam("user")
	password := getDbParam("password")

	return Connect(host, port, user, password, name)
}

func getDbParam(key string) string {
	return viper.GetString("application.database." + key)
}