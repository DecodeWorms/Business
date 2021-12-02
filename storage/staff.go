package storage

import (
	"business/types"

	"golang.org/x/crypto/bcrypt"
)

type Staff struct {
	cli Client
}

func NewStaff(cl Client) Staff {
	return Staff{
		cli: cl,
	}
}

func (stf Staff) AutoMigrate(data types.Product) error {
	return stf.cli.Client.AutoMigrate(&data)
}

func (stf Staff) SignUp(data types.Staff) error {
	pass, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	resPass := string(pass)
	cr := types.Staff{
		FullName: data.FullName,
		Password: resPass,
		Gender:   data.Gender,
		Unit:     data.Unit,
	}
	return stf.cli.Client.Create(&cr).Error
}

func (stf Staff) SignIn(data types.Staff) (types.Staff, error) {
	var s types.Staff
	return s, stf.cli.Client.Select("password").Where("full_name = ?", data.FullName).First(&s).Error

}

func (stf Staff) Save(data types.Product, name string) (types.Product, error) {
	res := data.Quantity * data.Price

	rd := types.Product{
		Name:         data.Name,
		AdminName:    name,
		SerialNumber: data.SerialNumber,
		Price:        data.Price,
		Quantity:     data.Quantity,
		Total:        res,
	}
	return rd, stf.cli.Client.Create(&rd).Error
}

func (stf Staff) Product(data types.Product) (types.Product, error) {
	var p types.Product
	return p, stf.cli.Client.First(&p, "serial_number = ?", data.SerialNumber).Error
}

func (stf Staff) Products() ([]types.Product, error) {
	var res []types.Product
	return res, stf.cli.Client.Find(&res).Error
}

func (stf Staff) Remove(srn string, data types.Product) error {
	return stf.cli.Client.Model(&types.Product{}).Where("serial_number = ?", srn).Delete(&data).Error
}
