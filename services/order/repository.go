package order

import (
	"log"

	"github.com/Ayobami6/pickitup/models"
	"github.com/Ayobami6/pickitup/services/order/dto"
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

func (o *OrderRepoImpl) GetOrders(userID uint)([]dto.OrderResponseDTO, error) {
	var orders []models.Order
    var orderResponse []dto.OrderResponseDTO

    res := o.db.Where(models.Order{UserID: userID}).Find(&orders)
    if res.Error!= nil {
        return orderResponse, res.Error
    }

    for _, order := range orders {
        orderResponse = append(orderResponse, dto.OrderResponseDTO{
			ID: order.ID,
			UserID: order.UserID, 
			RiderID: order.RiderID, 
			Status: string(order.Status), 
			CreatedAt: order.CreatedAt.String(),
			Charge: order.Charge,
            Item: order.Item,
			RefID: order.RefID,
			PaymentStatus: string(order.PaymentStatus),
			Acknowledge: order.Acknowledge,
			PickUpAddress: order.PickUpAddress,
		})
    }

    return orderResponse, nil
}

func (o *OrderRepoImpl) UpdateDeliveryStatus(orderId uint, status models.StatusType) error {
	return o.db.Model(&models.Order{}).Where("id =?", orderId).Update("status", status).Error
}

func (o *OrderRepoImpl) UpdateAcknowledgeStatus(orderID uint) error {
	return o.db.Model(&models.Order{}).Where("id =?", orderID).Update("acknowledge", true).Error
}

func (o *OrderRepoImpl) GetOrderByID(orderID uint) (dto.OrderResponseDTO, error) {
	var order models.Order
    var orderResponse dto.OrderResponseDTO

    res := o.db.Where("id =?", orderID).First(&order)
    if res.Error!= nil {
        return orderResponse, res.Error
    }

    orderResponse = dto.OrderResponseDTO{
        ID: order.ID,
        UserID: order.UserID,
        RiderID: order.RiderID,
        Status: string(order.Status),
        CreatedAt: order.CreatedAt.String(),
        Charge: order.Charge,
        Item: order.Item,
        RefID: order.RefID,
        PaymentStatus: string(order.PaymentStatus),
        Acknowledge: order.Acknowledge,
        PickUpAddress: order.PickUpAddress,
    }

    return orderResponse, nil
}