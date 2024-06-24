package controllers

import (
	"gotestbackend/database"
	"gotestbackend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var newUser models.User
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	// newUser.Password = string(hashedPassword)
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Simulate credit addition
	newUser.Credit = 1000.0
	// Save user to database (assuming db is initialized in main.go)
	database.DB.Create(&newUser)

	c.JSON(http.StatusCreated, newUser)
}

func GetAllUser(c *gin.Context) {
	var user models.User
	if err := database.DB.Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "All record not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUser retrieves the logged-in user's details
// @Summary      Get user by ID
// @Description  Get details of a user by their ID
// @Tags 		 users
// @Accept       json
// @Produce      json
// @Param        id      path    string  true  "User ID"
// @Success      200     {object} models.User
// @Failure      400     {object} models.ErrorResponse
// @Failure      404     {object} models.ErrorResponse
// @Router       /user/GetUserByID/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "User not found1111"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update user fields
	database.DB.Model(&user).Updates(updatedUser)

	c.JSON(http.StatusOK, user)
}

func DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Delete user
	database.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// LoginPayload is used to bind login request body
type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserPayload is used to bind update request body
type UpdateUserPayload struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Password      string `json:"password"`
	AccountNumber string `json:"account_number"`
}

// Login handles user login
func Login(c *gin.Context) {
	var payload LoginPayload
	var user models.User

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("username = ?", payload.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID})
}

// GetUser retrieves the logged-in user's details
// @Summary      Get by id
// @Description  get user profile
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id      path    string  true  "User ID"
// @Success      200     {object} models.User
// @Failure      400     {object} models.ErrorResponse
// @Failure      404     {object} models.ErrorResponse
// @Router       /user/profile [get]
func GetUserProfile(c *gin.Context) {
	var user models.User
	idparam := c.Param("user_id")

	//Convert idparam to uint
	userID, err := strconv.ParseUint(idparam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid user ID"})
		return
	}

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates the logged-in user's details
func UpdateUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	var payload UpdateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}
	if payload.AccountNumber != "" {
		user.AccountNumber = payload.AccountNumber
	}
	if payload.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	database.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

// api/accounting.go
// TransferCredit transfers credit from one user to another
func TransferCredit(c *gin.Context) {
	idparam := c.Param("user_id")

	//Convert idparam to uint
	userID, err := strconv.ParseUint(idparam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid user ID"})
		return
	}

	// Parse request body
	var transferRequest struct {
		ID uint `json:"id"`
		//SenderAccount   string  `json:"sender_account"`
		ReceiverAccount string  `json:"receiver_account"`
		Amount          float64 `json:"amount"`
	}
	transferRequest.ID = uint(userID)
	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implement transfer logic
	// Check if sender and receiver IDs are valid
	sender, err := GetDataUser(transferRequest.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender not found"})
		return
	}

	receiver, err := GetDataUserByAccount(transferRequest.ReceiverAccount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
		return
	}

	// Update sender and receiver credits in database
	// Validate if sender has enough credit
	if sender.Credit < transferRequest.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient credit"})
		return
	}

	// Perform credit transfer
	// db.Model(&sender).Update("credit", sender.Credit - amount)
	sender.Credit -= transferRequest.Amount
	// db.Model(&receiver).Update("credit", receiver.Credit + amount)
	receiver.Credit += transferRequest.Amount

	// Update sender and receiver in database
	err = database.DB.Save(&sender).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender"})
		return
	}

	err = database.DB.Save(&receiver).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receiver"})
		return
	}

	// Record transaction
	transaction := models.Transaction{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Amount:     transferRequest.Amount,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = database.DB.Create(&transaction).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}
func GetDataUser(user_id uint) (models.User, error) {
	var user models.User
	if err := database.DB.First(&user, user_id).Error; err != nil {
		return user, nil
	}
	return user, nil
}
func GetDataUserByAccount(account_number string) (models.User, error) {
	var user models.User
	if err := database.DB.First(&user, account_number).Error; err != nil {
		return user, nil
	}
	return user, nil
}

// TransferListRequest defines the query parameters for transfer list API
type TransferListRequest struct {
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
}

// GetTransferList retrieves the list of credit transfer history with optional filters
func GetTransferList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	var req TransferListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}
	// Prepare query conditions based on filters
	var transfers []models.Transaction
	db := database.DB
	query := db.Model(&models.Transaction{})
	query = query.Where("sender_id = ? OR receiver_id = ?", userID, userID)
	// Apply date filters if provided
	if req.StartDate != nil && req.EndDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	} else if req.StartDate != nil {
		query = query.Where("created_at >= ?", req.StartDate)
	} else if req.EndDate != nil {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	// Fetch transfers
	err := query.Find(&transfers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transfer history"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transfers": transfers})
}
