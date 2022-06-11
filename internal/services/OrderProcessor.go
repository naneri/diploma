package services

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/naneri/diploma/cmd/gophermart/controllers/dto"
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/user"
	"log"
	"time"
)

func ProcessOrders(UserRepo *user.DBRepository, ItemRepo *item.DBRepository, AccrualSystemAddress string) {
	unprocessedItems, queryErr := ItemRepo.GetUnprocessedItems()

	if queryErr != nil {
		log.Println(queryErr.Error())
		return
	}

	for _, dbItem := range unprocessedItems {
		var dtoOrder dto.Order
		time.Sleep(time.Second)
		client := resty.New()

		_, requestErr := client.R().SetResult(&dtoOrder).
			Get(fmt.Sprintf("%s/api/orders/%d", AccrualSystemAddress, dbItem.OrderID))

		if requestErr != nil {
			log.Println(requestErr.Error())
			return
		}

		err := ItemRepo.UpdateItemStatusAndAccrualByOrderID(dtoOrder.Order, dtoOrder.Status, dtoOrder.Accrual)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		updateBalanceErr := UserRepo.UpdateUserBalance(dbItem.UserID, dtoOrder.Accrual)

		if updateBalanceErr != nil {
			log.Println(updateBalanceErr.Error())
			continue
		}
	}
}
