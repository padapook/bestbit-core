package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/padapook/bestbit-core/internal/wallet/service"
	"github.com/shopspring/decimal"
)

type WalletController interface {
	GetWallets(c *gin.Context)
	GetWalletByCurrency(c *gin.Context)
	Deposit(c *gin.Context)
	Withdraw(c *gin.Context)
	Transfer(c *gin.Context)
}

type walletController struct {
	walletService service.WalletService
}

func NewWalletController(walletService service.WalletService) WalletController {
	return &walletController{walletService: walletService}
}

func (ctrl *walletController) GetWallets(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wallets, err := ctrl.walletService.GetUserWallets(accountID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wallets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": wallets,
	})
}

func (ctrl *walletController) GetWalletByCurrency(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currency := c.Param("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency is required"})
		return
	}

	wallet, err := ctrl.walletService.GetWalletBalance(accountID.(string), currency)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": wallet,
	})
}

type DepositRequest struct {
	Currency string          `json:"currency" binding:"required"`
	Amount   decimal.Decimal `json:"amount" binding:"required"`
}

func (ctrl *walletController) Deposit(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: amount and currency are required"})
		return
	}

	refID := uuid.New().String()

	wallet, err := ctrl.walletService.DepositMoney(accountID.(string), req.Currency, req.Amount, refID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Deposit successful",
		"data":    wallet,
	})
}

type WithdrawRequest struct {
	Currency string          `json:"currency" binding:"required"`
	Amount   decimal.Decimal `json:"amount" binding:"required"`
}

type TransferRequest struct {
	ToUserID string          `json:"to_user_id" binding:"required"`
	Currency string          `json:"currency" binding:"required"`
	Amount   decimal.Decimal `json:"amount" binding:"required"`
}

func (ctrl *walletController) Withdraw(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: amount and currency are required"})
		return
	}

	refID := uuid.New().String()

	wallet, err := ctrl.walletService.WithdrawMoney(accountID.(string), req.Currency, req.Amount, refID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Withdraw successful",
		"data":    wallet,
	})
}

func (ctrl *walletController) Transfer(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	refID := uuid.New().String()

	err := ctrl.walletService.TransferMoney(accountID.(string), req.ToUserID, req.Currency, req.Amount, refID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer successful",
	})
}
