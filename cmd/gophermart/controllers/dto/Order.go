package dto

type Order struct {
	Order   uint `json:"order,string"`
	Status  string
	Accrual float64
}
