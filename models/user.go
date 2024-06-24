package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID            uint    `json:"id" gorm:"primary_key"`
	Username      string  `json:"username"`
	Password      string  `json:"-"`
	AccountNumber string  `json:"account_number"`
	Credit        float64 `json:"credit"`
}
