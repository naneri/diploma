package controllers

import (
	"encoding/json"
	"errors"
	"github.com/naneri/diploma/cmd/gophermart/config"
	"github.com/naneri/diploma/cmd/gophermart/controllers/dto"
	"github.com/naneri/diploma/cmd/gophermart/httpservices"
	"github.com/naneri/diploma/cmd/gophermart/middleware"
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/user"
	"log"
	"net/http"
	"strconv"
	"time"
)

type OrderController struct {
	ItemRepo *item.DBRepository
	UserRepo *user.DBRepository
	Config   *config.Config
}

func (c OrderController) Add(w http.ResponseWriter, r *http.Request) {
	orderID, err := httpservices.GetOrderIDFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	loggedUser, userSearchErr := c.findUserFromRequest(r)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	order, orderFound, orderSearchErr := c.ItemRepo.GetItemByOrderID(orderID)

	if orderSearchErr != nil {
		http.Error(w, "Errors searching for the order", http.StatusInternalServerError)
		return
	}

	if orderFound {
		if order.UserID == loggedUser.ID {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	_, storeErr := c.ItemRepo.StoreItem(loggedUser.ID, orderID, item.StatusNew, 0)

	if storeErr != nil {
		http.Error(w, "error storing the order", http.StatusBadRequest)
		log.Println("error storing the order: " + storeErr.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (c OrderController) List(w http.ResponseWriter, r *http.Request) {

	loggedUser, userSearchErr := c.findUserFromRequest(r)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	items, err := c.ItemRepo.GetUserItems(loggedUser.ID)

	if err != nil {
		http.Error(w, "Errors searching user orders", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	outputItems := make([]dto.Item, 0)
	for _, dbItem := range items {
		outputItems = append(outputItems, dto.Item{
			Number:     strconv.FormatUint(uint64(dbItem.OrderID), 10),
			Status:     dbItem.Status,
			Accrual:    dbItem.Bonus,
			UploadedAt: time.Time{},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if len(outputItems) == 0 {
		return
	}

	encodeErr := json.NewEncoder(w).Encode(outputItems)

	if encodeErr != nil {
		http.Error(w, "error generating response", http.StatusInternalServerError)
		return
	}
}

func (c OrderController) findUserFromRequest(r *http.Request) (user.User, error) {
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
