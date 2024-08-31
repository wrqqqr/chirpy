package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type User struct {
	Email    string `json:"email"`
	Id       int    `json:"id"`
	Password []byte `json:"password"`
}

var ErrNotExist = errors.New("resource does not exist")

func NewDB(path string) (*DB, error) {

	newDB := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := newDB.ensureDB()

	return newDB, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Id:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	users, err := db.GetUsers()

	for _, v := range users {
		if v.Email == email {
			return User{}, fmt.Errorf("user already exist %v", err)
		}
	}

	fmt.Println(password, email)

	id := len(dbStructure.Users) + 1
	cryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("bad password converting")
		return User{}, err
	}

	user := User{
		Id:       id,
		Email:    email,
		Password: cryptedPass,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(data.Chirps))

	for _, v := range data.Chirps {
		chirps = append(chirps, v)
	}

	return chirps, nil
}

func (db *DB) GetUsers() ([]User, error) {
	data, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(data.Users))

	for _, v := range data.Users {
		users = append(users, v)
	}

	return users, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbs, err
	}

	err = json.Unmarshal(data, &dbs)

	if err != nil {
		return dbs, err
	}

	return dbs, nil

}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0666)

	if err != nil {
		return errors.New("something wrong")
	}
	return nil

}
