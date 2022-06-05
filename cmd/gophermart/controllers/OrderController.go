package controllers

import (
	"github.com/naneri/diploma/cmd/gophermart/middleware"
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/services"
	"github.com/naneri/diploma/internal/user"
	"io"
	"net/http"
	"strconv"
)

type OrderController struct {
	ItemRepo *item.DbRepository
	UserRepo *user.DbRepository
}

func (c OrderController) Add(w http.ResponseWriter, r *http.Request) {
	orderId, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value(middleware.UserID(middleware.UserIDContextKey)).(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	loggedUser, userSearchErr := c.UserRepo.FindUserById(userID)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	stringOrderId := string(orderId)

	intOrderId, parseErr := strconv.Atoi(stringOrderId)
	if parseErr != nil {
		http.Error(w, "wrong input format", http.StatusBadRequest)
		return
	}

	if !services.ValidLuhn(intOrderId) {
		http.Error(w, "Incorrect order id", http.StatusBadRequest)
		return
	}

	order, orderFound, orderSearchErr := c.ItemRepo.GetItemByOrderId(uint(intOrderId))

	if orderSearchErr != nil {
		http.Error(w, "Errors searching for the order", http.StatusInternalServerError)
		return
	}

	if orderFound {
		if order.UserId == loggedUser.ID {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
}

func (c OrderController) List(w http.ResponseWriter, r *http.Request) {

}
