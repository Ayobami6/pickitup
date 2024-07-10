package models

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"gorm.io/gorm"
)

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
	Verified bool `json:"verified" gorm:"default:false"`
	Rider      Rider      `gorm:"foreignKey:UserID"`
}

type Rider struct {
	ID uint  `json:"id" gorm:"primaryKey"`
	RiderID    string `json:"rider_id" gorm:"uniqueIndex;size:8"`
	BikeNumber string `json:"bike_number"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserID uint  `json:"user_id"`
	NextOfKinName string `json:"next_of_kin_name"`
	NextOfKinPhone string `json:"next_of_kin_phone"`
	DriverLicenseNumber string `json:"driver_license_number"`
	NextOfKinAddress string `json:"next_of_kin_address"`
	Address string `json:"address"`
	SuccessfulRides int64 `json:"successful_rides"`
	Rating float64  `json:"rating"`
	Level string  `json:"level"`
	CurrentLocation string `json:"current_location"`
	Reviews  []Review `json:"reviews" gorm:"foreignKey:RiderID"`
	CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
}

func (u *Rider) BeforeCreate(tx *gorm.DB) (err error) {
    u.RiderID, err = generateRandomID(8)
    if err != nil {
        return err
    }
	u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    return nil
}

func (r *Rider) BeforeUpdate(tx *gorm.DB) (err error) {
    r.UpdatedAt = time.Now()
    return nil
}



type Review struct {
	gorm.Model
	RiderID uint `json:"rider_id"`
	Rating float64 `json:"rating"`
	Comment string `json:"comment"`
}

func generateRandomID(n int) (string, error) {
    bytes := make([]byte, n)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return strings.ToUpper(hex.EncodeToString(bytes)[:n]), nil
}

