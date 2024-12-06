package controllers

import (
	"fmt"
	"net/http"
	"sales-app/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VoucherController struct {
	DB *gorm.DB
}

func NewVoucherController(db *gorm.DB) *VoucherController {
	return &VoucherController{DB: db}
}

func (vc *VoucherController) GetVouchers(c *gin.Context) {
	id := c.DefaultQuery("id", "")

	if id == "" {
		var vouchers []types.Voucher
		vc.DB.Preload("Brand").Find(&vouchers)
		c.JSON(http.StatusOK, vouchers)
	} else {
		var voucher types.Voucher
		if err := vc.DB.Preload("Brand").First(&voucher, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
			return
		}
		c.JSON(http.StatusOK, voucher)
	}
}

func (vc *VoucherController) GetVouchersByBrand(c *gin.Context) {
	id := c.DefaultQuery("id", "")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required."})
	} else {
		var vouchers []types.Voucher
		vc.DB.Preload("Brand").Where("brand_id = ?", id).Find(&vouchers)
		c.JSON(http.StatusOK, vouchers)
	}
}

func (vc *VoucherController) CreateVoucher(c *gin.Context) {
	var voucher types.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	var validationErrors []string

	if voucher.Name == "" {
		validationErrors = append(validationErrors, "Name is required.")
	}

	var brand types.Brand

	if voucher.BrandID == 0 {
		validationErrors = append(validationErrors, "Brand ID is required.")
	} else if err := vc.DB.First(&brand, voucher.BrandID).Error; err != nil {
		validationErrors = append(validationErrors, fmt.Sprintf("Brand with ID %d not found.", voucher.BrandID))
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	vc.DB.Create(&voucher)
	vc.DB.Preload("Brand").First(&voucher, voucher.ID)
	c.JSON(http.StatusOK, voucher)
}
