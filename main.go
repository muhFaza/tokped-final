package main

import (
	"net/http"
	"tokped-final/config"
	"tokped-final/handler"
	"tokped-final/helper"
	"tokped-final/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	h := &handler.Handler{DB: db}

	r := gin.Default()

	// uses bcrypt for password hashing
	r.POST("/users/register", h.Register)
	r.POST("/users/login", h.Login)

	r.PATCH("/users/topup", AuthMiddleware(h), h.TopUp)

	r.Run(":8080")
}

func AuthMiddleware(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		user := &model.User{}
		result := h.DB.Where("email = ?", claims.Email).First(user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func initDB() *gorm.DB {
	dbConfig := config.GetDBConfig()
	dsn := dbConfig.GetDBURL()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = db.AutoMigrate(&model.User{}, &model.Category{}, &model.Product{}, &model.TransactionHistory{})
	if err != nil {
		panic("Failed to migrate database!")
	}
	return db
}
