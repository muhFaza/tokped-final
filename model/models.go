package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null" validate:"required,email"`
	Fullname string `gorm:"not null;column:full_name" validate:"required"`
	Password string `gorm:"not null" validate:"required,min=6"`
	Role     string `gorm:"not null;default:customer"`
	Balance  int    `gorm:"not null;default:0"`
}

type Category struct {
	gorm.Model
	Type              string `gorm:"unique;not null" validate:"required"`
	SoldProductAmount int    `gorm:"column:sold_product_amount;default:0"`
}

type Product struct {
	gorm.Model
	Title      string `gorm:"not null" validate:"required"`
	Price      int    `gorm:"not null" validate:"required,min=5"`
	Stock      int    `gorm:"not null" validate:"required,min=0,max=50000000"`
	CategoryID int
	Category   Category `gorm:"foreignKey:CategoryID"`
}

type TransactionHistory struct {
	gorm.Model
	ProductID  int     `gorm:"not null"`
	Product    Product `gorm:"foreignKey:ProductID"`
	UserID     int     `gorm:"not null"`
	User       User    `gorm:"foreignKey:UserID"`
	Quantity   int     `gorm:"not null" validate:"required"`
	TotalPrice int     `gorm:"not null;column:total_price" validate:"required"`
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
