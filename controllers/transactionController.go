package controllers

import (
	"net/http"
	"sales-app/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{DB: db}
}

func (tc *TransactionController) GetTransactions(c *gin.Context) {
	id := c.DefaultQuery("transactionId", "")

	if id == "" {
		var transactions []types.Transaction
		tc.DB.Preload("Customer").Preload("Redemptions").Find(&transactions)
		c.JSON(http.StatusOK, transactions)
	} else {
		var transaction types.Transaction
		if err := tc.DB.Preload("Customer").Preload("Redemptions").First(&transaction, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusOK, transaction)
	}
}

func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var transaction types.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	var validationErrors []string

	if transaction.CustomerID == 0 {
		validationErrors = append(validationErrors, "Customer ID is required.")
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	var customer types.Customer
	if err := tc.DB.First(&customer, transaction.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found."})
		return
	}

	tc.DB.Preload("Customer").Create(&transaction)
	c.JSON(http.StatusOK, transaction)
}

func (tc *TransactionController) CreateRedemption(c *gin.Context) {
	var redemption types.Redemption
	if err := c.ShouldBindJSON(&redemption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	var validationErrors []string

	if redemption.TransactionID == 0 {
		validationErrors = append(validationErrors, "Transaction ID is required.")
	}

	if redemption.VoucherID == 0 {
		validationErrors = append(validationErrors, "Voucher ID is required.")
	}

	if redemption.Quantity == 0 {
		validationErrors = append(validationErrors, "Please enter the number of vouchers you would like to redeem.")
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	var transaction types.Transaction
	if err := tc.DB.First(&transaction, redemption.TransactionID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction not found."})
		return
	}

	var voucher types.Voucher
	if err := tc.DB.First(&voucher, redemption.VoucherID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voucher not found."})
		return
	}

	redemption.Points = voucher.CostInPoints * redemption.Quantity
	transaction.TotalPoints = transaction.TotalPoints + redemption.Points

	tc.DB.Create(&redemption)
	tc.DB.Save(&transaction)
	tc.DB.Preload("Transaction.Customer").Preload("Voucher.Brand").First(&redemption)
	c.JSON(http.StatusOK, redemption)
}
