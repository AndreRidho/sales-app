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

func TestGetBrands(t *testing.T) {
	gormDB, mock, err := test.SetupMockDB()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	brandController := controllers.NewBrandController(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `brands`")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(1, "Brand A", "Description of Brand A").
			AddRow(2, "Brand B", "Description of Brand B"))

	req, _ := http.NewRequest("GET", "/brands", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	brandController.GetBrands(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []types.Brand
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Brand A", response[0].Name)
	assert.Equal(t, "Brand B", response[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}

func TestCreateBrand(t *testing.T) {
	gormDB, mock, err := test.SetupMockDB()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	brandController := controllers.NewBrandController(gormDB)

	brand := types.Brand{
		Name:        "Brand X",
		Description: "Description of Brand X",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `brands`").
		WithArgs(brand.Name, brand.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body, _ := json.Marshal(brand)
	req, _ := http.NewRequest("POST", "/brands", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	brandController.CreateBrand(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response types.Brand
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, brand.Name, response.Name)
	assert.Equal(t, brand.Description, response.Description)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}

func TestCreateBrand_InvalidInput(t *testing.T) {
	gormDB, mock, err := test.SetupMockDB()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	brandController := controllers.NewBrandController(gormDB)

	brand := types.Brand{
		Name:        "",
		Description: "Description of Brand Y",
	}

	body, _ := json.Marshal(brand)
	req, _ := http.NewRequest("POST", "/brands", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	brandController.CreateBrand(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response["error"], "Name is required.")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}
