package httpServices

import (
	"errors"
	"github.com/naneri/diploma/internal/services"
	"io"
	"net/http"
	"strconv"
)

func GetOrderIdFromRequest(r *http.Request) (uint, error) {
	orderId, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		return 0, err
	}

	stringOrderId := string(orderId)

	return ParseOrderId(stringOrderId)
}

func ParseOrderId(stringOrderId string) (uint, error) {
	intOrderId, parseErr := strconv.Atoi(stringOrderId)
	if parseErr != nil {
		return 0, errors.New("wrong input format")
	}

	if !services.ValidLuhn(intOrderId) {
		return 0, errors.New("incorrect order id")
	}

	return uint(intOrderId), nil
}
