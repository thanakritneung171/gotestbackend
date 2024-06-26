package models

// @description User represents the entity of a user with basic information like username, personal details, account number, and credit balance.
type User struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// @description The account number associated with the user.
	AccountNumber string  `json:"account_number"`
	Credit        float64 `json:"credit"`
}
