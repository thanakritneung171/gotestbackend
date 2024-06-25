package models

type User struct {
	ID            uint    `json:"id" gorm:"primary_key"`
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	AccountNumber string  `json:"account_number"`
	Credit        float64 `json:"credit"`
}
