package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/naneri/diploma/cmd/gophermart/controllers/dto"
	"github.com/naneri/diploma/cmd/gophermart/httpservices"
	"github.com/naneri/diploma/cmd/gophermart/middleware"
	"github.com/naneri/diploma/internal/services"
	"github.com/naneri/diploma/internal/user"
	"github.com/naneri/diploma/internal/withdrawal"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type BalanceController struct {
	UserRepo       *user.DBRepository
	WithdrawalRepo *withdrawal.DBRepository
	DBConnection   *gorm.DB
}

func (c BalanceController) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	loggedUser, userSearchErr := c.findUserFromRequest(r)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	balance := dto.Balance{
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

	var WithdrawData dto.Withdraw

	if decodeErr := json.NewDecoder(r.Body).Decode(&WithdrawData); decodeErr != nil {
		http.Error(w, "please check that all fields are sent", http.StatusBadRequest)
		fmt.Println("error decoding json.")
		return
	}

	uintOrderID, parseErr := httpservices.ParseOrderID(WithdrawData.Order)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	withDrawErr := services.PerformWithdraw(c.DBConnection, c.UserRepo, loggedUser.ID, WithdrawData.Sum, uintOrderID)

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

func (c BalanceController) ListWithdrawals(w http.ResponseWriter, r *http.Request) {
	loggedUser, userSearchErr := c.findUserFromRequest(r)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	withdrawals, err := c.WithdrawalRepo.ListUserWithdrawals(loggedUser.ID)
	if err != nil {
		http.Error(w, "error getting withdrawals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	outputWithdrawals := make([]dto.OutputWithdrawal, 0)
	for _, dbWithdrawal := range withdrawals {
		outputWithdrawals = append(outputWithdrawals, dto.OutputWithdrawal{
			Order:       dbWithdrawal.OrderID,
			Sum:         dbWithdrawal.Sum,
			ProcessedAt: dbWithdrawal.CreatedAt,
		})
	}

	if len(outputWithdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(outputWithdrawals)

	if encodeErr != nil {
		http.Error(w, "error generating response", http.StatusInternalServerError)
		return
	}
}

func (c BalanceController) findUserFromRequest(r *http.Request) (user.User, error) {
	userID, ok := r.Context().Value(middleware.UserID(middleware.UserIDContextKey)).(uint32)

	if !ok {
		return user.User{}, errors.New("UserID not set in Cookie")
	}

	loggedUser, userSearchErr := c.UserRepo.FindUserByID(userID)
	if userSearchErr != nil {
		return user.User{}, userSearchErr
	}

	return loggedUser, nil
}
