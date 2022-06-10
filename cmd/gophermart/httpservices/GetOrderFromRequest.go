package httpservices

import (
	"errors"
	"github.com/naneri/diploma/internal/services"
	"io"
	"net/http"
	"strconv"
)

func GetOrderIDFromRequest(r *http.Request) (uint, error) {
	orderID, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		return 0, err
	}

	stringOrderID := string(orderID)

	return ParseOrderID(stringOrderID)
}

func ParseOrderID(stringOrderID string) (uint, error) {
	intOrderID, parseErr := strconv.Atoi(stringOrderID)
	if parseErr != nil {
		return 0, errors.New("wrong input format")
	}

	if !services.ValidLuhn(intOrderID) {
		return 0, errors.New("incorrect order id")
	}

	return uint(intOrderID), nil
}
