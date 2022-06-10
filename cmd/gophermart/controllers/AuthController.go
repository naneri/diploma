package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/naneri/diploma/cmd/gophermart/config"
	"github.com/naneri/diploma/cmd/gophermart/controllers/dto"
	"github.com/naneri/diploma/cmd/gophermart/httpservices"
	"github.com/naneri/diploma/internal/services"
	"github.com/naneri/diploma/internal/user"
	"log"
	"net/http"
)

type AuthController struct {
	UserRepo *user.DBRepository
	Config   *config.Config
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var authData dto.AuthData

	if decodeErr := json.NewDecoder(r.Body).Decode(&authData); decodeErr != nil {
		http.Error(w, "please check that all fields are sent", http.StatusBadRequest)
		fmt.Println("error decoding json.")
		return
	}

	checkErr := authData.CheckForEmptyFields()

	if checkErr != nil {
		http.Error(w, "please fill all required fields", http.StatusInternalServerError)
		fmt.Println("error checking the fields are ")
		return
	}

	_, userExists, searchErr := c.UserRepo.Find(authData.Login)

	if searchErr != nil {
		http.Error(w, searchErr.Error(), http.StatusInternalServerError)
		return
	}

	if userExists {
		http.Error(w, fmt.Sprintf("User with login %s already exists.", authData.Login), http.StatusConflict)
		return
	}

	newUser, saveErr := c.UserRepo.Save(authData.Login, authData.Password)

	if saveErr != nil {
		http.Error(w, "error storing the user", http.StatusInternalServerError)
		log.Println("error storing the user:" + saveErr.Error())
		return
	}

	authCookie := httpservices.GenerateUserCookie(newUser.ID, []byte(c.Config.SecretKey))
	http.SetCookie(w, &authCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var authData dto.AuthData

	if decodeErr := json.NewDecoder(r.Body).Decode(&authData); decodeErr != nil {
		http.Error(w, "please check that all fields are sent", http.StatusBadRequest)
		fmt.Println("error decoding json.")
		return
	}

	checkErr := authData.CheckForEmptyFields()

	if checkErr != nil {
		http.Error(w, "please fill all required fields", http.StatusInternalServerError)
		fmt.Println("error checking the fields are ")
		return
	}

	dbUser, userExists, searchErr := c.UserRepo.Find(authData.Login)

	if searchErr != nil || !userExists {
		http.Error(w, searchErr.Error(), http.StatusInternalServerError)
		return
	}

	if !services.CheckUserPassword(authData.Password, dbUser) {
		http.Error(w, "login or password incorrect", http.StatusUnauthorized)
		log.Printf("wrong login data for user %s \n", dbUser.Login)
		return
	}

	authCookie := httpservices.GenerateUserCookie(dbUser.ID, []byte(c.Config.SecretKey))
	http.SetCookie(w, &authCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
