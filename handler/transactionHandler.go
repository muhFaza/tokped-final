package handler

import (
	"net/http"
	"tokped-final/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (h *Handler) CreateTransaction(c *gin.Context) {
	transaction := &model.TransactionHistory{}
	if err := c.ShouldBindJSON(transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*model.User)
	transaction.UserID = int(user.ID)

	// check if product exists
	var product model.Product
	if err := h.DB.First(&product, transaction.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product does not exist"})
		return
	}

	transaction.TotalPrice = product.Price * transaction.Quantity

	validate := validator.New()
	if err := validate.Struct(transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user has enough balance
	if user.Balance < transaction.TotalPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough balance"})
		return
	}

	// check if product is in stock
	if product.Stock < transaction.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is out of stock"})
		return
	}

	// create transaction
	if err := h.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update product stock
	product.Stock = product.Stock - transaction.Quantity
	if err := h.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update user balance
	user.Balance = user.Balance - transaction.TotalPrice
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "You have successfully purchased the product",
		"transaction_bill": gin.H{
			"total_price":   transaction.TotalPrice,
			"quantity":      transaction.Quantity,
			"product_title": product.Title,
		}})
}
