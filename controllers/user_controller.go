package controllers

import (
	"net/http"

	"gotestbackend/database"
	"gotestbackend/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var newUser models.User
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
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
