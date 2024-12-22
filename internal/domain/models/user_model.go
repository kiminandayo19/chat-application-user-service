package models

import (
	"time"
)

type Mst_users struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string `json:"username" gorm:"index;unique"`
	Email       string `json:"email" gorm:"unique"`
	Phonenumber string `json:"phoneNumber" gorm:"unique"`
	Password    string `json:"password"`
	IsOnline    bool   `json:"isOnline"`
	IsDeleted   bool   `json:"isDeleted"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
