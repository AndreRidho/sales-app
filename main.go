package main

import (
	"log"
	"sales-app/controllers"
	"sales-app/types"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDatabase() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/sales_app?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&types.Brand{}, &types.Voucher{}, &types.Customer{}, &types.Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	return db
}

func main() {
	db := initDatabase()

	brandController := controllers.NewBrandController(db)
	customerController := controllers.NewCustomerController(db)
	voucherController := controllers.NewVoucherController(db)
	transactionController := controllers.NewTransactionController(db)

	r := gin.Default()

	r.GET("/brands", brandController.GetBrands)
	r.POST("/brands", brandController.CreateBrand)

	r.GET("/customers", customerController.GetCustomers)
	r.POST("/customers", customerController.CreateCustomer)

	r.GET("/vouchers", voucherController.GetVouchers)
	r.GET("/vouchers/brand", voucherController.GetVouchersByBrand)
	r.POST("/vouchers", voucherController.CreateVoucher)

	r.GET("/transactions/redemption", transactionController.GetTransactions)
	r.POST("/transactions/redemption", transactionController.CreateRedemption)
	r.POST("/transactions", transactionController.CreateTransaction)

	r.Run(":8080")
}
