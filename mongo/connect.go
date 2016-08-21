package mongo

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

var _currentMongoSession *mgo.Session
var _currentMongoDatabase *mgo.Database

// GetMongoSession get the current mongo session
func GetMongoSession() *mgo.Session {
	return _currentMongoSession
}

// GetMongoDatabase get the current mongo database
func GetMongoDatabase() *mgo.Database {
	return _currentMongoDatabase
}

// ConnectMongo connect to the mongo database
func ConnectMongo(host, port, user, password, name string) (session *mgo.Session, err error) {
	mongooseConnectionChain := "mongodb://" +
		map[bool]string{true: user, false: ""}[len(user) > 0] +
		map[bool]string{true: ":" + password, false: ""}[len(password) > 0] +
		map[bool]string{true: "@", false: ""}[len(user) > 0] +
		host +
		map[bool]string{true: ":" + port, false: ""}[len(port) > 0] +
		map[bool]string{true: "/" + name, false: ""}[len(name) > 0]

	_currentMongoSession, err = mgo.Dial(mongooseConnectionChain)

	if len(name) > 0 {
		_currentMongoDatabase = _currentMongoSession.DB(name)
	}

	fmt.Println(_currentMongoSession.DatabaseNames())

	if err != nil {
		fmt.Println(mongooseConnectionChain)
		panic(err.Error())
	}

	return
}
