package service

import (
	"context"
	"errors"
	// "log"

	"github.com/google/uuid"

	accountModel "github.com/padapook/bestbit-core/internal/account/model"
	"github.com/padapook/bestbit-core/internal/account/repository"
	"github.com/padapook/bestbit-core/internal/utils"
	"github.com/padapook/bestbit-core/internal/utils/auth"
	"gorm.io/gorm"

	walletModel "github.com/padapook/bestbit-core/internal/wallet/model"

	"github.com/shopspring/decimal"
)

type UserService interface {
	Register(ctx context.Context, user accountModel.User) (*accountModel.User, error)
	GetByUsername(username string) (*accountModel.User, error)
	Login(username, password string) (*accountModel.User, error)
	LoginByShareToken(token string) (*accountModel.User, error)
}

type userService struct {
	repo repository.UserRepository
	db   *gorm.DB
}

func NewUserService(repo repository.UserRepository, db *gorm.DB) UserService {
	return &userService{repo: repo, db: db}
}

func (s *userService) Register(ctx context.Context, user accountModel.User) (*accountModel.User, error) {
	// pwhashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	pwhashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.Password = pwhashed

	user.AccountId = uuid.New().String()

	var createdUser *accountModel.User

	err = s.db.Transaction(func(tx *gorm.DB) error {
		usernameExist, _ := s.repo.GetByUsername(user.Username)
		if usernameExist != nil && usernameExist.ID != 0 {
			return errors.New("username already exists")
		}

		if err := s.repo.CreateUser(tx, &user); err != nil {
			return err
		}

		newWallet := &walletModel.Wallet{
			UserID:   user.ID,
			Currency: "THB",
			Balance:  decimal.NewFromInt(0),
		}

		if err := tx.Create(newWallet).Error; err != nil {
			return err
		}

		createdUser = &user
		return nil
	})

	return createdUser, err
}

func (s *userService) GetByUsername(username string) (*accountModel.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *userService) Login(username, password string) (*accountModel.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	match, err := utils.ComparePasswordAndHash(password, user.Password)
	if err != nil || !match {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *userService) LoginByShareToken(token string) (*accountModel.User, error) {
	// log.Println("'token",token)
	claims, err := auth.ValidateShareToken(token)
	// log.Println("'claims",claims)
	if err != nil {
		return nil, errors.New("invalid or expired share token")
	}

	user, err := s.repo.GetByUsername(claims.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
