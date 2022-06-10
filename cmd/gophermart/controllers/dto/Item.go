package dto

import (
	"encoding/json"
	"time"
)

type Item struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (i Item) MarshalJSON() ([]byte, error) {
	type ItemAlias Item

	aliasValue := struct {
		ItemAlias
		// переопределяем поле внутри анонимной структуры
		UploadedAt string `json:"uploaded_at"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		ItemAlias: ItemAlias(i),
		// задаём значение для переопределённого поля
		UploadedAt: i.UploadedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}
