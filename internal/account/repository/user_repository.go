package repository

import (
	"github.com/padapook/bestbit-core/internal/account/model"
    "gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(tx *gorm.DB, user *model.User) error
	GetByUsername(username string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) CreateUser(tx *gorm.DB, user *model.User) error {
    return tx.Create(user).Error
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	
	err := r.db.Where("username = ?", username).First(&user).Error

	return &user, err
}