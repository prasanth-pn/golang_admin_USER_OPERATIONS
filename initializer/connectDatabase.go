package initializer

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {

	// dns-data sorce name

	dsn := os.Getenv("DB")

	// Establish connections to the database and returning if any error

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect ")

	}
	fmt.Println("connected successfully")

}
