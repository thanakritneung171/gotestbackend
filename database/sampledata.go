package database

import (
	"gotestbackend/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func InsertSampleUser() {
	var count int64
	// Check if the users table is empty
	if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
		log.Fatalf("Failed to count users: %v", err)
	}

	if count > 0 {
		log.Println("Users table already has data. Skipping sample data insertion.")
		return
	}
	users := []models.User{
		{Username: "user1", Password: hashPassword("password1"), FirstName: "John", LastName: "Doe", AccountNumber: "1111111111"},
		{Username: "user2", Password: hashPassword("password2"), FirstName: "Jane", LastName: "Doe", AccountNumber: "2222222222"},
		{Username: "user3", Password: hashPassword("password3"), FirstName: "Alice", LastName: "Smith", AccountNumber: "3333333333"},
		{Username: "user4", Password: hashPassword("password4"), FirstName: "Bob", LastName: "Brown", AccountNumber: "4444444444"},
		{Username: "user5", Password: hashPassword("password5"), FirstName: "Charlie", LastName: "Davis", AccountNumber: "5555555555"},
		{Username: "user6", Password: hashPassword("password6"), FirstName: "David", LastName: "Evans", AccountNumber: "6666666666"},
		{Username: "user7", Password: hashPassword("password7"), FirstName: "Ella", LastName: "Green", AccountNumber: "7777777777"},
		{Username: "user8", Password: hashPassword("password8"), FirstName: "Frank", LastName: "Harris", AccountNumber: "8888888888"},
		{Username: "user9", Password: hashPassword("password9"), FirstName: "Grace", LastName: "Johnson", AccountNumber: "9999999999"},
		{Username: "user10", Password: hashPassword("password10"), FirstName: "Henry", LastName: "Lee", AccountNumber: "1010101010"},
	}

	for _, user := range users {
		// Set default credit value
		user.Credit = 1000
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Could not insert user %s: %v", user.Username, err)
		}
	}

}

func InsertSampleTransaction() {
	var count int64
	// Check if the users table is empty
	if err := DB.Model(&models.Transaction{}).Count(&count).Error; err != nil {
		log.Fatalf("Failed to count Transaction: %v", err)
	}

	if count > 0 {
		log.Println("Transaction table already has data. Skipping sample data insertion.")
		return
	}
	// Insert sample transactions
	transactions := []models.Transaction{
		{SenderID: 1, ReceiverID: 2, Amount: 100},
		{SenderID: 2, ReceiverID: 3, Amount: 200},
		{SenderID: 3, ReceiverID: 4, Amount: 150},
		{SenderID: 4, ReceiverID: 5, Amount: 120},
		{SenderID: 5, ReceiverID: 6, Amount: 80},
		{SenderID: 6, ReceiverID: 7, Amount: 60},
		{SenderID: 7, ReceiverID: 8, Amount: 90},
		{SenderID: 8, ReceiverID: 9, Amount: 110},
		{SenderID: 9, ReceiverID: 10, Amount: 130},
		{SenderID: 10, ReceiverID: 1, Amount: 50},
	}

	for _, transaction := range transactions {
		if err := DB.Create(&transaction).Error; err != nil {
			log.Printf("Could not insert transaction from user %d to user %d: %v", transaction.SenderID, transaction.ReceiverID, err)
		}
	}
}

// Hash password using bcrypt
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}
