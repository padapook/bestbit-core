package database

import (
	accountModel "github.com/padapook/bestbit-core/internal/account/model"
	orderModel "github.com/padapook/bestbit-core/internal/order/model"
	tradeModel "github.com/padapook/bestbit-core/internal/trade/model"
	walletModel "github.com/padapook/bestbit-core/internal/wallet/model"

	"gorm.io/gorm"
	"log"
)

func AutoMigrate(db *gorm.DB) error {
	// log.Println("เข้า migrate")

	err := db.AutoMigrate(
		// account
		&accountModel.User{},

		// wallet
		&walletModel.Transaction{},
		&walletModel.Wallet{},

		// order
		&orderModel.Order{},

		// trade
		&tradeModel.Trade{},
	)

	if err != nil {
		log.Println("'พัง")
		return err
	}

	log.Println("migrate success")
	return nil
}
