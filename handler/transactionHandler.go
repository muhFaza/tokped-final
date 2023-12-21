package handler

import (
	"net/http"
	"tokped-final/model"

	"github.com/gin-gonic/gin"
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

	// update sold product amount in category
	var category model.Category
	if err := h.DB.First(&category, product.CategoryID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	category.SoldProductAmount += transaction.Quantity
	if err := h.DB.Save(&category).Error; err != nil {
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

type GetMyTransactionsResponse struct {
	ID         uint
	ProductID  int
	Product    model.Product
	UserID     int
	Quantity   int
	TotalPrice int
}

func (h *Handler) GetMyTransactions(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	var transactions []model.TransactionHistory
	result := h.DB.Preload("Product").Where("user_id = ?", user.ID).Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var response []GetMyTransactionsResponse
	for _, transaction := range transactions {
		response = append(response, GetMyTransactionsResponse{
			ID:         transaction.ID,
			ProductID:  transaction.ProductID,
			Product:    transaction.Product,
			Quantity:   transaction.Quantity,
			TotalPrice: transaction.TotalPrice,
			UserID:     transaction.UserID,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetTransactions(c *gin.Context) {
	var transactions []model.TransactionHistory
	result := h.DB.Preload("Product").Preload("User").Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
