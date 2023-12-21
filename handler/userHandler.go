package handler

import (
	"net/http"
	"tokped-final/model"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) Register(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Role = "customer"
	user.Balance = 0

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.DB.Create(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "full_name": user.Fullname, "email": user.Email, "balance": user.Balance, "created_at": user.CreatedAt})
}
