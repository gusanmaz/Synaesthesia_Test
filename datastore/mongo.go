package datastore

import "labix.org/v2/mgo"

type Mongo struct {
	url     string
	Session *mgo.Session
	connErr error
}

func (db *Mongo) New(url string) {
	db.url = url
	db.Session, db.connErr = mgo.Dial("127.0.0.1:27017")
	// TODO May instead failsafe connection problem
	if db.connErr != nil {
		panic(db.connErr)
	}
}

func (db *Mongo) Close() {
	db.Session.Close()
}
