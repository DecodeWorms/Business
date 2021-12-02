package types

import "gorm.io/gorm"

type Staff struct {
	gorm.Model
	FullName string `json:"full_name" gorm:"fullname" required:"required,min=10"`
	Password string `json:"password" gorm:"password" required:"passwd,required,min=10"`
	Gender   string `json:"gender" gorm:"gender" required:"required,min=8"`
	Unit     string `json:"unit" gorm:"unit" required:"required"`
}
