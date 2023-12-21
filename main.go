package main

import (
	"tokped-final/config"
	"tokped-final/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = initDB()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong!",
		})
	})

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
