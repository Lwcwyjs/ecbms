package handler

import (
	"net/http"
	"strconv"

	"ecbms/internal/model"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := model.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := model.CreateUser(req.Username, req.Password, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	model.AddSystemLog("user", "Created user: "+req.Username, "info")

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := model.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	model.AddSystemLog("user", "Deleted user ID: "+idParam, "info")

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
