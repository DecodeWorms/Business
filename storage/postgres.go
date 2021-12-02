package storage

import (
	"business/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	Client *gorm.DB
}

func NewClient(cfg config.Config) Client {
	log.Println("Connecting to DB..")
	dsn := fmt.Sprintf("host=%s dbname=%s port=%s user=%s", cfg.DatabaseHost, cfg.DatabaseName, cfg.DatabasePort, cfg.DatabaseUsername)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Unable to connect to DB..")
	}
	log.Println("Connected...")

	return Client{
		Client: db,
	}

}
