package controllers

import (
	"net/http"
	"sales-app/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandController struct {
	DB *gorm.DB
}

func NewBrandController(db *gorm.DB) *BrandController {
	return &BrandController{DB: db}
}

func (bc *BrandController) GetBrands(c *gin.Context) {
	var brands []types.Brand
	bc.DB.Preload("Vouchers").Find(&brands)
	c.JSON(http.StatusOK, brands)
}

func (bc *BrandController) CreateBrand(c *gin.Context) {
	var brand types.Brand
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	var validationErrors []string

	if brand.Name == "" {
		validationErrors = append(validationErrors, "Name is required.")
	}

	if brand.Description == "" {
		validationErrors = append(validationErrors, "Description is required.")
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	bc.DB.Create(&brand)
	c.JSON(http.StatusOK, brand)
}
