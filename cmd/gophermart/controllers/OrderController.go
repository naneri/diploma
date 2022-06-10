package controllers

import (
	"encoding/json"
	"errors"
	"github.com/naneri/diploma/cmd/gophermart/config"
	"github.com/naneri/diploma/cmd/gophermart/controllers/Dto"
	"github.com/naneri/diploma/cmd/gophermart/httpServices"
	"github.com/naneri/diploma/cmd/gophermart/middleware"
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/user"
	"log"
	"net/http"
	"strconv"
	"time"
)

type OrderController struct {
	ItemRepo *item.DbRepository
	UserRepo *user.DbRepository
	Config   *config.Config
}

func (c OrderController) Add(w http.ResponseWriter, r *http.Request) {
	orderId, err := httpServices.GetOrderIdFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	loggedUser, userSearchErr := c.findUserFromRequest(r)
	if userSearchErr != nil {
		http.Error(w, "error finding the user", http.StatusInternalServerError)
		return
	}

	order, orderFound, orderSearchErr := c.ItemRepo.GetItemByOrderId(orderId)

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

	_, storeErr := c.ItemRepo.StoreItem(loggedUser.ID, orderId, item.StatusNew, 0)

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

	outputItems := make([]Dto.Item, 0)
	for _, dbItem := range items {
		outputItems = append(outputItems, Dto.Item{
			Number:     strconv.FormatUint(uint64(dbItem.OrderId), 10),
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
		return user.User{}, errors.New("UserId not set in Cookie")
	}

	loggedUser, userSearchErr := c.UserRepo.FindUserById(userID)
	if userSearchErr != nil {
		return user.User{}, userSearchErr
	}

	return loggedUser, nil
}
