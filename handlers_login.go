package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	user := User{}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	users, err := cfg.DB.GetUsers()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users")
	}

	for _, v := range users {
		if v.Email == params.Email {
			err := bcrypt.CompareHashAndPassword([]byte(v.Password), []byte(params.Password))
			fmt.Println(v)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "unauthorized ")
				return
			} else {
				user = User{
					Email: v.Email,
					Id:    v.Id,
				}
				respondWithJSON(w, http.StatusOK, user)
				return
			}

		}
	}

	respondWithJSON(w, http.StatusOK, user)
}
