package controllers

import (
	"encoding/json"
	"errors"
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
	"time"
)

type OrderController struct {
	ItemRepo *item.DbRepository
	UserRepo *user.DbRepository
	Config   *config.Config
}

func (c OrderController) Add(w http.ResponseWriter, r *http.Request) {
	orderId, err := c.getOrderIdFromRequest(r)

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

	// _________
	// _________   NEED TO Move this to a dedicated process that will review the
	// _________
	client := resty.New()

	resp, requestErr := client.R().Get(fmt.Sprintf("%s/api/orders/%d", c.Config.AccrualAddress, orderId))
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
		updateErr := c.UserRepo.UpdateUserBalance(loggedUser.ID, dtoOrder.Accrual)
		if updateErr != nil {
			log.Println("error updating the balance: ", updateErr)
		}
	}
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

func (c OrderController) getOrderIdFromRequest(r *http.Request) (uint, error) {
	orderId, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		return 0, err
	}

	stringOrderId := string(orderId)

	intOrderId, parseErr := strconv.Atoi(stringOrderId)
	if parseErr != nil {
		return 0, errors.New("wrong input format")
	}

	if !services.ValidLuhn(intOrderId) {
		return 0, errors.New("incorrect order id")
	}

	return uint(intOrderId), nil
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
