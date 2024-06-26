package controllers

import (
	"gotestbackend/database"
	"gotestbackend/middlewares"
	"gotestbackend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//	@Summary		Register a new user
//	@Description	Registers a new user with initial credit
//	@Tags			Auth , CRUD
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User data"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	map[string]string	"message"
//	@Failure		401		{object}	map[string]string	"message"
//	@Router			/user/register [post]
func Register(c *gin.Context) {
	var newUser models.User
	//fmt.Println("passhash :", string(hashedPassword))
	if err := c.ShouldBindJSON(&newUser); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusNotFound, gin.H{"Message": "Invalid input"})
		return
	}
	// Validate AccountNumber
	if !isValidAccountNumber(newUser.AccountNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Account Number must be a 10-digit number"})
		return
	}
	var userexists models.User
	database.DB.Where("username = ? ", newUser.Username).First(&userexists)
	if userexists.ID > 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "User exist"})
		return
	}
	database.DB.Where("account_number = ? ", newUser.AccountNumber).First(&userexists)
	if userexists.ID > 0 {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Account Number exist"})
		return
	}
	//fmt.Println("pass :", newUser.Password)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashedPassword)
	//fmt.Println("pass hashedPassword:", newUser.Password)
	// Simulate credit addition
	newUser.Credit = 1000.0
	// Save user to database
	database.DB.Create(&newUser)

	c.JSON(http.StatusCreated, newUser)
}

// isValidAccountNumber validates if the account number is a 10-digit number
func isValidAccountNumber(accountNumber string) bool {
	if len(accountNumber) != 10 {
		return false
	}
	for _, r := range accountNumber {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

//	@Summary		Get All User
//	@Description	Get details all user
//	@Tags			CRUD
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.User
//	@Failure		404	{object}	map[string]string	"message"
//	@Router			/userAll [get]
func GetAllUser(c *gin.Context) {
	var user []models.User
	if err := database.DB.Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "All record not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUser retrieves the logged-in user's details
//
//	@Summary		Get user by ID
//	@Description	Get details of a user by their ID
//	@Tags			CRUD
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/user/GetUserByID/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

//	@Summary		Update User By ID
//	@Description	Update a user
//	@Tags			CRUD
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//
//	@Param			id		path		string		true	"User ID"
//
//	@Param			user	body		models.User	true	"User data"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	map[string]string	"message"
//	@Failure		404		{object}	map[string]string	"message"
//	@Router			/user/UpdateUserByID/{id} [put]
func UpdateUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Update user fields
	database.DB.Model(&user).Updates(updatedUser)

	c.JSON(http.StatusOK, user)
}

//	@Summary		Update User By ID
//	@Description	Update a  user
//	@Tags			CRUD
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"User ID"
//	@Success		201	{object}	map[string]string	"message"
//	@Failure		400	{object}	map[string]string	"message"
//	@Router			/user/DeleteUserByID/{id} [delete]
func DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
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

// Login godoc
//
//	@Summary		login
//	@Description	Authenticates a user and returns a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		LoginPayload		true	"Login payload"
//	@Success		200		{object}	map[string]string	"token"
//	@Failure		400		{object}	map[string]string	"message"
//	@Failure		401		{object}	map[string]string	"message"
//	@Failure		500		{object}	map[string]string	"message"
//	@Router			/user/login [post]
func Login(c *gin.Context) {
	var payload LoginPayload
	var user models.User

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	database.DB.Where("username =?", payload.Username).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user dose not exists"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "pass login"})
	}
	//fmt.Println("User ID = ", user.ID)
	token, err := middlewares.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token}) //token
}

// GetUser retrieves the logged-in user's details
//
//	@Summary		getUser
//	@Description	GetUserProfile by id token
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/user/me [get]
func GetUser(c *gin.Context) {
	var user models.User
	idparam, exists := c.Get("user_id")
	//fmt.Println("user_id ", idparam)
	if !exists {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "User ID not found in context"})
		return
	}
	// Convert userID to the appropriate type
	userID, ok := idparam.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "User ID type assertion failed"})
		return
	}
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser updates the logged-in user's details
//
//	@Summary		updateUser
//	@Description	UpdateUser by id token
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UpdateUserPayload	false	"UserPayload data"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	map[string]string	"message"
//	@Failure		401		{object}	map[string]string	"message"
//	@Failure		404		{object}	map[string]string	"message"
//	@Router			/user/me [patch]
func UpdateUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in"})
		return
	}
	var payload UpdateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}
	if payload.AccountNumber != "" {
		// Validate AccountNumber
		if !isValidAccountNumber(payload.AccountNumber) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account Number must be a 10-digit number"})
			return
		}
		var userexists models.User
		database.DB.Where("account_number = ? ", payload.AccountNumber).First(&userexists)
		if userexists.ID > 0 {
			c.JSON(http.StatusNotFound, gin.H{"Message": "Account Number exist"})
			return
		}
		user.AccountNumber = payload.AccountNumber
	}
	if payload.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}
	database.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

type transferRequest struct {
	//ID uint `json:"id"`
	//SenderAccount   string  `json:"sender_account"`
	ReceiverAccount string  `json:"receiver_account"`
	Amount          float64 `json:"amount"`
}

// TransferCredit transfers credit from one user to another
//
//	@Summary		transfer
//	@Description	TransferCredit transfers credit from one user to another
//	@Tags			accounting
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			transferRequest	body		transferRequest	true	"transferRequest data"
//	@Success		200				{object}	models.Transaction
//	@Failure		400				{object}	map[string]string	"message"
//	@Failure		404				{object}	map[string]string	"message"
//	@Failure		500				{object}	map[string]string	"message"
//	@Router			/accounting/transfer [post]
func Transfer(c *gin.Context) {
	idparam, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in"})
		return
	}
	// Convert userID to the appropriate type
	userID, ok := idparam.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID"})
		return
	}
	// Parse request body
	var transferRequest transferRequest
	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Implement transfer logic
	// Check if sender and receiver IDs are valid
	sender, err := GetDataUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Sender not found"})
		return
	}
	receiver, err := GetDataUserByAccount(transferRequest.ReceiverAccount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Receiver not found"})
		return
	}
	// Update sender and receiver credits in database
	// Validate if sender has enough credit
	if sender.Credit < transferRequest.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient credit"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update sender"})
		return
	}
	err = database.DB.Save(&receiver).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update receiver"})
		return
	}
	// Record transaction
	transaction := models.Transaction{
		SenderID:          sender.ID,
		SenderRemaining:   sender.Credit,
		ReceiverID:        receiver.ID,
		ReceiverRemaining: receiver.Credit,
		Amount:            transferRequest.Amount,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = database.DB.Create(&transaction).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to record transaction"})
		return
	}
	c.JSON(http.StatusOK, transaction)
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
	if err := database.DB.Where("account_number = ?", account_number).First(&user).Error; err != nil {
		return user, nil
	}
	//fmt.Println("GetDataUserByAccount account_number : ", account_number)
	//fmt.Println("GetDataUserByAccount user : ", user.ID)
	return user, nil
}

// TransferListRequest defines the query parameters for transfer list API
type TransferListRequest struct {
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
}

// GetTransferList retrieves the list of credit transfer history with optional filters
//
//	@Summary		getTransferList
//	@Description	GetTransferList retrieves the list of credit transfer history with optional filters
//	@Tags			accounting
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			TransferListRequest	body		TransferListRequest	false	"date: '2024-06-25'"
//	@Success		200					{object}	models.Transaction
//	@Failure		400					{object}	map[string]string	"message"
//	@Failure		401					{object}	map[string]string	"message"
//	@Failure		500					{object}	map[string]string	"message"
//	@Router			/accounting/transfer-list [get]
func GetTransferList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}
	var req TransferListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters"})
		return
	}
	// Prepare query conditions based on filters
	//fmt.Println("StartDate =", string(req.StartDate))
	//fmt.Println("EndDate =", string(req.EndDate))
	layout := "2006-01-02"
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		parsedStartDate, err := time.Parse(layout, req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date"})
			return
		}
		startDate = &parsedStartDate
	}
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse(layout, req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date"})
			return
		}
		endDate = &parsedEndDate
		*endDate = endDate.AddDate(0, 0, 1)
	}
	//fmt.Println("StartDate =", startDate)
	//fmt.Println("EndDate =", endDate)
	var transfers []models.Transaction
	db := database.DB
	query := db.Model(&models.Transaction{})
	query = query.Where("sender_id = ? OR receiver_id = ?", userID, userID)
	// Apply date filters if provided
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startDate, *endDate)
		//fmt.Println("BETWEEN")
	} else if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
		//fmt.Println("StartDate != nil")
	} else if endDate != nil {
		query = query.Where("created_at < ?", *endDate)
		//fmt.Println("EndDate != nil")
	}
	// Fetch transfers
	err := query.Find(&transfers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transfer history"})
		return
	}
	c.JSON(http.StatusOK, transfers)
}
