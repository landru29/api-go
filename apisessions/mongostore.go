package apisessions

import (
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	mgo "gopkg.in/mgo.v2"
)

// MongoStore is the structure for session storage
type MongoStore interface {
	sessions.Store
}

// NewMongoStore create a session store
func NewMongoStore(c *mgo.Collection, maxAge int, ensureTTL bool, keyPairs ...[]byte) MongoStore {
	return &mongoStore{mongostore.NewMongoStore(c, maxAge, ensureTTL, keyPairs...)}
}

type mongoStore struct {
	*mongostore.MongoStore
}

func (c *mongoStore) Options(options sessions.Options) {
	c.MongoStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
