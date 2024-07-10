package user

import (
	"log"

	"github.com/Ayobami6/pickitup/models"
	"gorm.io/gorm"
)

type userRepoImpl struct {
	db *gorm.DB
}


func NewUserRepoImpl(db *gorm.DB) *userRepoImpl {
	err := db.AutoMigrate(&models.User{})
	if err!= nil {
        log.Fatal(err)
    }
	return &userRepoImpl{db: db}
}


func (r *userRepoImpl) CreateUser(user *models.User) error {
	res := r.db.Create(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *userRepoImpl) GetUserByEmail(email string) (*models.User, error) {
	result := &models.User{}
	err := r.db.Where("email =?", email).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepoImpl) GetUserByID(id int) (*models.User, error) {
	result := &models.User{}
    err := r.db.First(&result, id).Error
    if err!= nil {
        return nil, err
    }
    return result, nil
}