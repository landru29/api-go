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

var currentSession *mgo.Session

func GetSession() (*mgo.Session) {
	return currentSession
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
	map[bool]string{true: ":" + port, false:""} [len(port)>0] +
	"/" + name;

	session, err =  mgo.Dial(mongooseConnectionChain)
	currentSession = session

	if err != nil {
		panic(err)
	}
	defer session.Close()

	return
}

func getDbParam(key string) string {
	return viper.GetString("application.database." + key)
}