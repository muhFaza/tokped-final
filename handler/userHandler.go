package handler

import (
	"fmt"
	"net/http"
	"tokped-final/helper"
	"tokped-final/model"

	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	result := h.DB.Create(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "full_name": user.Fullname, "email": user.Email, "balance": user.Balance, "created_at": user.CreatedAt})
}

func (h *Handler) Login(c *gin.Context) {
	loginRequest := &model.LoginRequest{}
	if err := c.ShouldBindJSON(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{}
	result := h.DB.Where("email = ?", loginRequest.Email).First(user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := helper.CreateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) TopUp(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	topUpRequest := &model.TopUpRequest{}
	if err := c.ShouldBindJSON(topUpRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(topUpRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Balance += topUpRequest.Balance
	result := h.DB.Save(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	res := fmt.Sprintf("Your balance has been successfully updated to %s", helper.FormatIDR(user.Balance))

	c.JSON(http.StatusOK, gin.H{"message": res})
}