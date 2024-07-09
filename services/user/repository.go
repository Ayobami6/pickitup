package user

import (
	"log"

	"gorm.io/gorm"
)

type userRepoImpl struct {
	db *gorm.DB
}


func NewUserRepoImpl(db *gorm.DB) *userRepoImpl {
	err := db.AutoMigrate(&User{})
	if err!= nil {
        log.Fatal(err)
    }
	return &userRepoImpl{db: db}
}


func (r *userRepoImpl) CreateUser(user *User) error {
	res := r.db.Create(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *userRepoImpl) GetUserByEmail(email string) (*User, error) {
	result := &User{}
	err := r.db.Where("email =?", email).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}