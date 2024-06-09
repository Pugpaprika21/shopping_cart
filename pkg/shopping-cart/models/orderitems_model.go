package models

import "gorm.io/gorm"

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Completed OrderStatus = "completed"
	Cancelled OrderStatus = "cancelled"
)

type OrderItems struct {
	gorm.Model
	ProductID       uint        `gorm:"product_id"`
	ProductQuantity uint        `gorm:"product_quantity"`
	UserID          uint        `gorm:"user_id"`
	UserIP          string      `gorm:"user_ip"`
	Status          OrderStatus `gorm:"status"` // ใช้ประเภท OrderStatus
}
