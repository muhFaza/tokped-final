package handler

import (
	"net/http"
	"time"
	"tokped-final/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	product := &model.Product{}
	if err := c.ShouldBindJSON(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if category exists
	var category model.Category
	if err := h.DB.First(&category, product.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category does not exist"})
		return
	}

	// Check if a soft deleted record with the same name exists
	var existingProduct model.Product
	if err := h.DB.Unscoped().Where("title = ?", product.Title).First(&existingProduct).Error; err == nil {
		// If a soft deleted record exists, undelete it
		h.DB.Model(&existingProduct).Update("deleted_at", nil)
		c.JSON(http.StatusOK, existingProduct)
		return
	}

	// If no soft deleted record exists, create a new one
	if err := h.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": product.ID, "title": product.Title, "price": product.Price, "stock": product.Stock, "category_id": product.CategoryID, "created_at": product.CreatedAt})
}

type ProductResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Title      string    `json:"title"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID int       `json:"category_Id"`
}

func (h *Handler) GetProducts(c *gin.Context) {
	var products []model.Product
	result := h.DB.Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var response []ProductResponse
	for _, product := range products {
		response = append(response, ProductResponse{
			ID:         product.ID,
			CreatedAt:  product.CreatedAt,
			Title:      product.Title,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryID: product.CategoryID,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	var newProduct model.Product
	var currProduct model.Product

	if err := h.DB.First(&newProduct, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product does not exist"})
		return
	}

	currProduct = newProduct

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newProduct.Title == "" {
		newProduct.Title = currProduct.Title
	}

	if newProduct.Price == 0 {
		newProduct.Price = currProduct.Price
	}

	if newProduct.Stock == 0 {
		newProduct.Stock = currProduct.Stock
	}

	if newProduct.CategoryID == 0 {
		// check if category id exists
		var category model.Category
		if err := h.DB.First(&category, currProduct.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category does not exist"})
			return
		}

		newProduct.CategoryID = currProduct.CategoryID
	}

	if err := h.DB.Save(&newProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": gin.H{"id": newProduct.ID, "title": newProduct.Title, "price": newProduct.Price, "stock": newProduct.Stock, "CategoryId": newProduct.CategoryID, "createdAt": newProduct.CreatedAt, "updatedAt": newProduct.UpdatedAt}})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Product
	result := h.DB.Where("id = ?", id).First(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if err := h.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product has been successfully deleted"})
}
