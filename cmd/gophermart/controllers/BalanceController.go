package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/naneri/diploma/cmd/gophermart/controllers/Dto"
	"github.com/naneri/diploma/cmd/gophermart/httpServices"
	"github.com/naneri/diploma/cmd/gophermart/middleware"
	"github.com/naneri/diploma/internal/services"
	"github.com/naneri/diploma/internal/user"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type BalanceController struct {
	UserRepo     *user.DbRepository
	DbConnection *gorm.DB
}

func (c BalanceController) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	loggedUser, userSearchErr := c.findUserFromRequest(r)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	balance := Dto.Balance{
		Current:   loggedUser.Balance,
		Withdrawn: loggedUser.WithDrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(balance)

	if encodeErr != nil {
		http.Error(w, "error generating response", http.StatusInternalServerError)
		return
	}
}

func (c BalanceController) RequestWithdraw(w http.ResponseWriter, r *http.Request) {
	loggedUser, userSearchErr := c.findUserFromRequest(r)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	var WithdrawData Dto.Withdraw

	if decodeErr := json.NewDecoder(r.Body).Decode(&WithdrawData); decodeErr != nil {
		http.Error(w, "please check that all fields are sent", http.StatusBadRequest)
		fmt.Println("error decoding json.")
		return
	}

	uintOrderId, parseErr := httpServices.ParseOrderId(WithdrawData.Order)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	withDrawErr := services.PerformWithdraw(c.DbConnection, c.UserRepo, loggedUser.ID, WithdrawData.Sum, uintOrderId)

	var lowBalanceErr services.NotEnoughBalanceError

	if errors.As(withDrawErr, &lowBalanceErr) {
		http.Error(w, "not enough balance", http.StatusPaymentRequired)
		return
	}

	if withDrawErr != nil {
		http.Error(w, "error processing the withdraw", http.StatusInternalServerError)
		log.Println(withDrawErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c BalanceController) findUserFromRequest(r *http.Request) (user.User, error) {
	userID, ok := r.Context().Value(middleware.UserID(middleware.UserIDContextKey)).(uint32)

	if !ok {
		return user.User{}, errors.New("UserId not set in Cookie")
	}

	loggedUser, userSearchErr := c.UserRepo.FindUserById(userID)
	if userSearchErr != nil {
		return user.User{}, userSearchErr
	}

	return loggedUser, nil
}
