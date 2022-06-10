package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/naneri/diploma/cmd/gophermart/config"
	"github.com/naneri/diploma/cmd/gophermart/controllers/Dto"
	"github.com/naneri/diploma/cmd/gophermart/middleware"
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/services"
	"github.com/naneri/diploma/internal/user"
	"io"
	"log"
	"net/http"
	"strconv"
)

type OrderController struct {
	ItemRepo *item.DbRepository
	UserRepo *user.DbRepository
	Config   *config.Config
}

func (c OrderController) Add(w http.ResponseWriter, r *http.Request) {
	orderId, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
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
		http.Error(w, "wrong input format", http.StatusUnprocessableEntity)
		return
	}

	if !services.ValidLuhn(intOrderId) {
		http.Error(w, "Incorrect order id", http.StatusUnprocessableEntity)
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

	_, storeErr := c.ItemRepo.StoreItem(userID, uint(intOrderId), item.StatusNew, 0)

	if storeErr != nil {
		http.Error(w, "error storing the order", http.StatusBadRequest)
		log.Println("error storing the order: " + storeErr.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)

	client := resty.New()

	resp, requestErr := client.R().Get(fmt.Sprintf("%s/api/orders/%d", c.Config.AccrualAddress, intOrderId))
	if requestErr != nil {
		http.Error(w, "Errors searching for the order in the accrual system", http.StatusInternalServerError)
		return
	}

	var dtoOrder Dto.Order

	if respDecode := json.NewDecoder(resp.RawBody()).Decode(&dtoOrder); respDecode != nil {
		http.Error(w, "error interacting with internal system", http.StatusBadRequest)
		log.Println("error decoding Order json from internal system: " + respDecode.Error())
		return
	}

	if dtoOrder.Accrual != 0 {
		updateErr := c.UserRepo.UpdateUserBalance(userID, dtoOrder.Accrual)
		if updateErr != nil {
			log.Println("error updating the balance: ", updateErr)
		}
	}
}

func (c OrderController) List(w http.ResponseWriter, r *http.Request) {

}
