package controllers

import (
	"net/http"
	"sales-app/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerController struct {
	DB *gorm.DB
}

func NewCustomerController(db *gorm.DB) *CustomerController {
	return &CustomerController{DB: db}
}

func (cc *CustomerController) GetCustomers(c *gin.Context) {
	var customers []types.Customer
	cc.DB.Find(&customers)
	c.JSON(http.StatusOK, customers)
}

func (cc *CustomerController) CreateCustomer(c *gin.Context) {
	var customer types.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	var validationErrors []string

	if customer.Name == "" {
		validationErrors = append(validationErrors, "Name is required.")
	}

	if customer.Email == "" {
		validationErrors = append(validationErrors, "E-mail is required.")
	} else {
		var existingCustomer types.Customer
		if err := cc.DB.Where("email = ?", customer.Email).First(&existingCustomer).Error; err == nil {
			validationErrors = append(validationErrors, "E-mail already exists.")
		}
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	cc.DB.Create(&customer)
	c.JSON(http.StatusOK, customer)
}
