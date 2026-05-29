package handler

import (
	"net/http"

	"ecbms/config"
	"ecbms/internal/model"
	"ecbms/pkg/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !model.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	cfg := config.Load()
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	model.AddSystemLog("auth", "User "+user.Username+" logged in", "info")

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
		"role":     role,
	})
}

func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, _ := c.Get("user_id")
	user, err := model.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	if !model.CheckPassword(user.Password, req.OldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	if err := model.UpdateUserPassword(user.ID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	model.AddSystemLog("user", "User "+user.Username+" changed password", "info")

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
