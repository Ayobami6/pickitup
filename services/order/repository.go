package order

import (
	"log"

	"github.com/Ayobami6/pickitup/models"
	"gorm.io/gorm"
)

type OrderRepoImpl struct {
	db *gorm.DB
}

func NewOrderRepoImpl(db *gorm.DB) *OrderRepoImpl{
	err := db.AutoMigrate(&models.Order{})
	if err != nil {
		log.Println(err)
	}
	return &OrderRepoImpl{db: db}
}

//  implement and interface model later 

func (o *OrderRepoImpl) CreateOrder(order *models.Order) error {
    return o.db.Create(order).Error
}