package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StatusType string

type PaymentStatus string

type RiderAvailabilityStatus string

const (
	Available RiderAvailabilityStatus = "Available"
    Unavailable RiderAvailabilityStatus = "Unavailable"
    OnBreak RiderAvailabilityStatus = "On Break"
	Busy RiderAvailabilityStatus = "Busy"
)

const (
	Paid PaymentStatus = "Paid"
	Unpaid PaymentStatus = "Unpaid"
)

const (
    Delivered StatusType = "Delivered"
    Canceled StatusType = "Canceled"
    Pending StatusType = "Pending"
	InDelivery StatusType = "InDelivery"
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserName string `json:"username"`
	Email string `json:"email required" gorm:"unique"`
	PhoneNumber string `json:"phone_number" gorm:"unique"`
	WalletBalance float64 `json:"wallet_balance"`
	AccountNumber string `json:"account_number" gorm:""`
	AccountName string `json:"account_name"`
	BankName string `json:"bank_name"`
	Password string `json:"password"`
	Verified bool `json:"verified" gorm:"default:false"`
	Rider      Rider      `gorm:"foreignKey:UserID"`
	Orders 	[]Order `json:"orders" gorm:"foreignKey:UserID"`
	CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
}

func (u *User) Debit(tx *gorm.DB, amount float64) error  {
	// Debit user wallet balance and save user again
	u.WalletBalance = u.WalletBalance - amount
    return tx.Save(u).Error

}

func (u *User) Credit(tx *gorm.DB, amount float64) error {
	// fund user wallet
	u.WalletBalance = u.WalletBalance + amount
	return tx.Save(u).Error
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    u.UpdatedAt = time.Now()
    return nil
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
	Orders 	[]Order `json:"orders" gorm:"foreignKey:RiderID"`
	CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
	// add minimum and maximum charge to this
	MinimumCharge float64 `json:"minimum_charge"`
	MaximumCharge float64 `json:"maximum_charge"`
	AvailabilityStatus RiderAvailabilityStatus `json:"availability_status" gorm:"default:Available"`
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

func (r *Rider)UpdateSuccessfulRides(tx *gorm.DB) error {
	r.SuccessfulRides++
    return tx.Save(r).Error
}

func (r *Rider) BeforeUpdate(tx *gorm.DB) (err error) {
    r.UpdatedAt = time.Now()
    return nil
}

type Order struct {
	ID uint `json:"id"`
	RiderID uint `json:"rider_id"`
	UserID uint `json:"user_id"`
	RefID string `json:"ref_id" gorm:"uniqueIndex;size:30"`
	Status StatusType `json:"status" gorm:"default:Pending"`
	Item string `json:"item"`
	Quantity int `json:"quantity"`
	Charge float64 `json:"price"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	PickUpAddress string `json:"pickup_address"`
	DropOffAddress string `json:"dropoff_address"`
	CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
	Acknowledge bool `json:"acknowledge" gorm:"default:false"`

}

func(u *Order) BeforeSave(tx *gorm.DB) error {
	u.Status = Pending
	if u.Status != Delivered && u.Status != Canceled && u.Status != Pending && u.Status != InDelivery {
		return fmt.Errorf("invalid status")
	}
	u.PaymentStatus = Unpaid
	if u.PaymentStatus != Paid && u.PaymentStatus != Unpaid {
		return fmt.Errorf("invalid payment status")
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	Id, err := generateRandomID(4)
	if err!= nil {
        return err
    }
	id := strings.ToUpper(Id)
	now := time.Now()
	year, month, day := now.Date()
    hour, minute, _ := now.Clock()
	o.RefID = fmt.Sprintf("PICK%s%d%d%d%d%d",id,year,month,day,hour,minute)
	o.CreatedAt = time.Now()
    o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) BeforeUpdate(tx *gorm.DB) (err error) {
    o.UpdatedAt = time.Now()
	if o.Acknowledge {
		fmt.Println("I got here")
		o.Status = InDelivery
	}
    return nil
}

type Review struct {
	ID		uint `json:"id"`
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

