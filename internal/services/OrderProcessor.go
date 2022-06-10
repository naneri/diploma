package services

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/naneri/diploma/cmd/gophermart/controllers/Dto"
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/user"
	"log"
	"time"
)

func ProcessOrders(UserRepo *user.DbRepository, ItemRepo *item.DbRepository, AccrualSystemAddress string) {
	unprocessedItems, queryErr := ItemRepo.GetUnprocessedItems()

	if queryErr != nil {
		log.Println(queryErr.Error())
		return
	}

	for _, dbItem := range unprocessedItems {
		var dtoOrder Dto.Order
		time.Sleep(time.Second)
		client := resty.New()

		_, requestErr := client.R().SetResult(&dtoOrder).
			Get(fmt.Sprintf("%s/api/orders/%d", AccrualSystemAddress, dbItem.OrderId))

		if requestErr != nil {
			log.Println(requestErr.Error())
			return
		}

		err := ItemRepo.UpdateItemStatusAndAccrualByOrderId(dtoOrder.Order, dtoOrder.Status, dtoOrder.Accrual)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		updateBalanceErr := UserRepo.UpdateUserBalance(dbItem.UserId, dtoOrder.Accrual)

		if updateBalanceErr != nil {
			log.Println(updateBalanceErr.Error())
			continue
		}
	}
}
