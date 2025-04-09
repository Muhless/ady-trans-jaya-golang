package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=MUHLESS717GG dbname=ady-trans-jaya port=5432 sslmode=disable Timezone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed connect to Database: %v", err)
	}
	fmt.Println("Succesfully connect to Database")
	DB = db
}
