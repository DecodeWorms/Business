package types

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name         string `json:"p_name" gorm:"fullname" required:"required"`
	AdminName    string `json:"admin_name" gorm:"admin_name"`
	SerialNumber string `json:"serial_number" gorm:"serialNumber" required:"required"`
	Price        int    `json:"price" gorm:"price" required:"required"`
	Quantity     int    `json:"quantity" gorm:"quantity" required:"required"`
	Total        int    `json:"total" gorm:"total" required:"required"`
}
