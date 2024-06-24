package models

type Transaction struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	SenderID   uint    `json:"sender_id"`
	ReceiverID uint    `json:"receiver_id"`
	Amount     float64 `json:"amount"`
}
