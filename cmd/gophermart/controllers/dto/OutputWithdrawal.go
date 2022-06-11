package dto

import (
	"encoding/json"
	"time"
)

type OutputWithdrawal struct {
	Order       uint      `json:"order,string"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (o OutputWithdrawal) MarshalJSON() ([]byte, error) {
	type WithdrawalAlias OutputWithdrawal

	aliasValue := struct {
		WithdrawalAlias
		// переопределяем поле внутри анонимной структуры
		ProcessedAt string `json:"processed_at"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		WithdrawalAlias: WithdrawalAlias(o),
		// задаём значение для переопределённого поля
		ProcessedAt: o.ProcessedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}
