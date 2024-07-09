package rider

import (
	"log"

	"gorm.io/gorm"
)

type riderRepositoryImpl struct {
	db *gorm.DB
}

func NewRiderRepositoryImpl(db *gorm.DB) *riderRepositoryImpl {
	err := db.AutoMigrate(&Rider{})
	if err!= nil {
        log.Fatal(err)
    }
    return &riderRepositoryImpl{db: db}
}

// overide the interface methods

func (r *riderRepositoryImpl) CreateRider(rider *Rider) error {
    return r.db.Create(rider).Error
}
