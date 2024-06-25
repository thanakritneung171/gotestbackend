package models

import (
	"time"
)

type Transaction struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	SenderID          uint      `json:"sender_id"`
	SenderRemaining   float64   `json:"ender_remaining "`
	ReceiverID        uint      `json:"receiver_id"`
	ReceiverRemaining float64   `json:"receiver_remaining "`
	Amount            float64   `json:"amount"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
