package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sales-app/controllers"
	"sales-app/test"
	"sales-app/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCustomers(t *testing.T) {
	gormDB, mock, err := test.SetupMockDB()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	customerController := controllers.NewCustomerController(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customers`")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "Customer A", "a@mail.com").
			AddRow(2, "Customer B", "b@mail.com"))

	req, _ := http.NewRequest("GET", "/customers", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	customerController.GetCustomers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []types.Customer
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Customer A", response[0].Name)
	assert.Equal(t, "Customer B", response[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}

func TestCreateCustomer(t *testing.T) {
	gormDB, mock, err := test.SetupMockDB()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	customerController := controllers.NewCustomerController(gormDB)

	customer := types.Customer{
		Name:  "Customer X",
		Email: "x@mail.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `customers`").
		WithArgs(customer.Name, customer.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	customerController.CreateCustomer(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response types.Customer
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, customer.Name, response.Name)
	assert.Equal(t, customer.Email, response.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}

func TestCreateCustomer_InvalidInput(t *testing.T) {
	gormDB, mock, err := test.SetupMockDB()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	customerController := controllers.NewCustomerController(gormDB)

	customer := types.Customer{
		Name:  "Customer X",
		Email: "x@mail.com",
	}

	body, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	customerController.CreateCustomer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response["error"], "Name is required.")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}
