package initializer

import "github.com/prasanth-pn/admin-login/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Admin{})
}
