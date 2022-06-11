package services

import "fmt"

type NotEnoughBalanceError struct {
	balance float64
	sum     float64
}

func (e NotEnoughBalanceError) Error() string {
	return fmt.Sprintf("not enough bonuses. Balance: %f, Requested Withdraw Amount: %f", e.balance, e.sum)
}
