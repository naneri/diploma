package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/naneri/diploma/cmd/gophermart/config"
	"github.com/naneri/diploma/cmd/gophermart/controllers/Dto"
	"github.com/naneri/diploma/cmd/gophermart/httpServices"
	"github.com/naneri/diploma/internal/user"
	"log"
	"net/http"
)

type AuthController struct {
	UserRepo *user.DbRepository
	Config   *config.Config
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var authData Dto.AuthData

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
		http.Error(w, fmt.Sprintf("User with login %s already exists.", authData.Login), http.StatusInternalServerError)
		return
	}

	newUser, saveErr := c.UserRepo.Save(authData.Login, authData.Password)

	if saveErr != nil {
		http.Error(w, "error storing the user", http.StatusInternalServerError)
		log.Println("error storing the user:" + saveErr.Error())
		return
	}

	authCookie := httpServices.GenerateUserCookie(newUser.ID, []byte(c.Config.SecretKey))
	http.SetCookie(w, &authCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c AuthController) Login(w http.ResponseWriter, r *http.Request) {

}
