package user

import (
	"errors"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

// User holds data for a single user
type User struct {
	ID   bson.ObjectId `json:"id" storm:"id"`
	Name string        `json:"name"`
	Role string        `json:"role"`
}

const (
	dbPath = "users.db"
)

var (
	// ErrRecordInvalid is the error showing a record is invalid
	ErrRecordInvalid = errors.New("record is invalid")
)

// All retrieves all users from DB
func All() ([]User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	users := []User{}
	err = db.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// One retrieves the user of given id
func One(id bson.ObjectId) (*User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	u := new(User)
	err = db.One("ID", id, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Delete removes the user of given id
func Delete(id bson.ObjectId) error {
	// must retrieve first
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	u := new(User)
	err = db.One("ID", id, u)
	if err != nil {
		return err
	}
	return db.DeleteStruct(u)
}

// Save creates a new record
func (u *User) Save() error {
	if err := u.validate(); err != nil {
		return err
	}
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Save(u)
}

// Validate checks the data inside a given user
func (u *User) validate() error {
	if u.Name == "" {
		return ErrRecordInvalid
	}
	return nil
}
