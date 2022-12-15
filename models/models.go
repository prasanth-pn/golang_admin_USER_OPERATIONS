package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string
	Password string
	Status   string
}
type Admin struct {
	gorm.Model
	Email    string
	Password string
}
