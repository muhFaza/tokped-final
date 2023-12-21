package main

import (
	"tokped-final/config"
	"tokped-final/handler"
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

	r.Run(":8080")
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
