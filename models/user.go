package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	FullName  string `json:"fullName" gorm:"not null"`
	UserName  string `json:"username" gorm:"unique; not null"`
	Password  string `json:"password" gorm:"not null"`
	Email     string `json:"email" gorm:"unique; not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
