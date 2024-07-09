package rider

import (
	"log"

	"github.com/Ayobami6/pickitup/models"
	"gorm.io/gorm"
)

type riderRepositoryImpl struct {
	db *gorm.DB
}

func NewRiderRepositoryImpl(db *gorm.DB) *riderRepositoryImpl {
	err := db.AutoMigrate(&models.Rider{}, &models.Review{})
	if err!= nil {
        log.Fatal(err)
    }
    return &riderRepositoryImpl{db: db}
}

// overide the interface methods

func (r *riderRepositoryImpl) CreateRider(rider *models.Rider) error {
    return r.db.Create(rider).Error
}
