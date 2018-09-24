package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
)


//Mongo configuration along with the collection name
type MongoConfig struct {
	Server     string
	Database   string
	Collection string
}

//create a connection with mongo and get its session
func (m *MongoConfig) Connect() (*mgo.Session, error) {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return session, nil
}


func (m *MongoConfig) getmgoDB(session *mgo.Session) (*mgo.Database, error) {

	var DB *mgo.Database

	DB = session.DB(m.Database)

	c := session.DB(m.Database).C(m.Collection)

	// Index
	index := mgo.Index{
		//This is hardcoded here but can be modified and get info as function parameter/ as part of MongoConfig
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return DB, nil
}

//session close
func Close(session *mgo.Session) {
	session.Close()
}

type UserInfo map[string]interface{}

// Find list in the collection
func (m *MongoConfig) FindAll() ([]UserInfo, error) {
	session, err := m.Connect()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer Close(session)
	DB, err := m.getmgoDB(session)
	var usersinfo []UserInfo
	err = DB.C(m.Collection).Find(bson.M{}).All(&usersinfo)

	return usersinfo, err
}

// Find a userinfo by unique id (key value in mgo.index)
func (m *MongoConfig) FindByUsername(username string) (UserInfo, error) {
	session, err := m.Connect()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer Close(session)
	DB, err := m.getmgoDB(session)

	var userinfo UserInfo
	err = DB.C(m.Collection).Find(bson.M{"name": username}).One(&userinfo)
	return userinfo, err
}

// Insert a new userinfo into Collection.
func (m *MongoConfig) Insert(userinfo UserInfo) error {
	session, err := m.Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer Close(session)
	DB, err := m.getmgoDB(session)

	err = DB.C(m.Collection).Insert(&userinfo)
	return err
}

// Delete  existing by unique key value in mgo.Index
func (m *MongoConfig) Delete(username string) error {
	session, err := m.Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer Close(session)
	DB, err := m.getmgoDB(session)

	userinfo, err := m.FindByUsername(username)
	if err != nil {
		return err
	}
	err = DB.C(m.Collection).Remove(&userinfo)
	return err
}

// Update existing
func (m *MongoConfig) Update(username string, update UserInfo) error {
	session, err := m.Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer Close(session)
	DB, err := m.getmgoDB(session)

	err = DB.C(m.Collection).Update(bson.M{"name": username}, update)
	return err
}

//Delete tags
func (m *MongoConfig) DeleteTags(username string, tags []string) error {
	userinfo, err := m.FindByUsername(username)
	if err != nil {
		return err
	}
	for _, value := range tags {
		delete(userinfo, value)
	}
	err = m.Update(username, userinfo)
	if err != nil {
		return err
	}
	return nil
}
//Update and also Insert
func (m *MongoConfig) Upsert(username string, update UserInfo) error {
	session, err := m.Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer Close(session)
	DB, err := m.getmgoDB(session)

	userinfo, err := m.FindByUsername(username)
	if err != nil {
		return err
	}
	for key, value := range update {
		userinfo[key] = value
	}

	_, err = DB.C(m.Collection).Upsert(bson.M{"name": username}, userinfo)
	return err
}
