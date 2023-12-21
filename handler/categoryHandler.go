package handler

import (
	"net/http"
	"tokped-final/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	category := &model.Category{}
	if err := c.ShouldBindJSON(category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check if a soft deleted record with the same type exists
	var existingCategory model.Category
	if err := h.DB.Unscoped().Where("type = ?", category.Type).First(&existingCategory).Error; err == nil {
			// If a soft deleted record exists, undelete it
			h.DB.Model(&existingCategory).Update("deleted_at", nil)
			c.JSON(http.StatusOK, existingCategory)
			return
	}

	validate := validator.New()
	if err := validate.Struct(category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If no soft deleted record exists, create a new one
	if err := h.DB.Create(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusCreated, gin.H{"id": category.ID, "type": category.Type, "sold_product_amount": category.SoldProductAmount, "created_at": category.CreatedAt})
}

func (h *Handler) GetCategories(c *gin.Context) {
	var categories []model.Category
    result := h.DB.Preload("Products").Find(&categories)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    var response []model.CategoryResponse
    for _, category := range categories {
        response = append(response, model.CategoryResponse{
            ID:                category.ID,
            CreatedAt:         category.CreatedAt,
            UpdatedAt:         category.UpdatedAt,
            Type:              category.Type,
            SoldProductAmount: category.SoldProductAmount,
						Products:          category.Products,
        })
    }

    c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category model.Category
	result := h.DB.Where("id = ?", id).First(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var newCategory model.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// replace the old category type with the new one
	category.Type = newCategory.Type

	result = h.DB.Model(&category).Updates(category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": category.ID, "type": category.Type, "sold_product_amount": category.SoldProductAmount, "created_at": category.CreatedAt})

}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category model.Category
	result := h.DB.Where("id = ?", id).First(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = h.DB.Delete(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category has been successfully deleted"})
}