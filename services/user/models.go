package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserName string `json:"username"`
	Email string `json:"email required" gorm:"unique"`
	PhoneNumber string `json:"phoneNumber" gorm:"unique"`
	WalletBalance string `json:"walletBalance"`
	AccountNumber string `json:"accountNumber" gorm:""`
	AccountName string `json:"accountName"`
	Password string `json:"password"`
}