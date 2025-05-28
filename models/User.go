package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"not null" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	Username  string `gorm:"not null; unique_index" json:"username"`
	Tasks     []Task `json:"tasks"`
}