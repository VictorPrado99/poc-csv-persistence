package csvpersistencecontroller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	c "github.com/VictorPrado99/poc-csv-persistence/internal/csv_persistence_controller"
	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
	"github.com/VictorPrado99/poc-csv-persistence/pkg/api"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSetupRouter(t *testing.T) {
	router := c.SetupRoutes()

	if router == nil {
		t.Fatalf("Failed to create router")
	}
}

func StartMockMySqlDb() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()

	dialec := mysql.New(mysql.Config{
		Conn: db,
	})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).WillReturnRows(mock.NewRows([]string{"version"}).AddRow(1))
	database.Instance, _ = gorm.Open(dialec, &gorm.Config{})

	return db, mock
}

func TestCreateOrdersSuccess(t *testing.T) {
	_, mock := StartMockMySqlDb()

	data := api.Order{
		Id:           4,
		Email:        "test@test.com",
		PhoneNumber:  "055 11 6583 2753",
		ParcelWeight: 54.74,
		Date:         "2022-05-31",
		Country:      "Brazil",
	}

	json, _ := json.Marshal(data)

	r := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(json))

	w := httptest.NewRecorder()

	mock.ExpectBegin()
	mock.ExpectCommit()

	c.CreateOrders(w, r)

	if w.Result().StatusCode != 201 {
		t.Fatalf("Didn't receive status 201, received status %d", w.Result().StatusCode)
	}

	w.Result().Body.Close()
}

func TestCreateOrdersFail(t *testing.T) {
	_, mock := StartMockMySqlDb()

	data := api.Order{
		Id:           4,
		Email:        "test@test.com",
		PhoneNumber:  "055 11 6583 2753",
		ParcelWeight: 54.74,
		Date:         "2022-05-31",
		Country:      "Brazil",
	}

	jsonData, _ := json.Marshal(data)

	r := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(jsonData))

	w := httptest.NewRecorder()

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("fail to save data"))

	c.CreateOrders(w, r)

	if w.Result().StatusCode != 208 {
		t.Fatalf("Didn't receive status 208, received status %d", w.Result().StatusCode)
	}

	w.Result().Body.Close()
}

// A test get without filter analysing the query return, not the best solution, but just to test a hypothesis. A better version of this method is at controller2_test.go
func TestGetWithoutFilters(t *testing.T) {
	_, mock := StartMockMySqlDb() // Generate a db mock instance

	orderType := reflect.TypeOf(api.Order{})
	var columns []string

	// Get columns of the entity
	for i := 0; i < orderType.NumField(); i++ {
		columns = append(columns, orderType.Field(i).Name)
	}

	// Generate mock for request and write
	r := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	// Query we are expecting, end what will define if the test is a success or not
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` ORDER BY id asc LIMIT 10")).WillReturnRows(mock.NewRows(columns))

	// Run the actual endpoint
	c.GetOrders(w, r)
}

func TestGetBadWeightLimit(t *testing.T) {
	// Generate mock for request and write
	r := httptest.NewRequest(http.MethodGet, "/orders?weightLimit=Foo", nil)
	w := httptest.NewRecorder()

	// Run the actual endpoint
	c.GetOrders(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("Status Code different from %d", http.StatusBadRequest)
	}
}

func TestHeadCount(t *testing.T) {
	_, mock := StartMockMySqlDb() // Generate a db mock instance

	// Generate mock for request and write
	r := httptest.NewRequest(http.MethodHead, "/orders/Brazil", nil)
	w := httptest.NewRecorder()

	// Query we are expecting, end what will define if the test is a success or not
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) as order_count, sum(parcel_weight) as total FROM `orders` WHERE country = ? ORDER BY `orders`.`id` LIMIT 1")).
		WithArgs("Brazil").
		WillReturnRows(mock.NewRows([]string{"order_count", "parcel_weight"}).AddRow(2, 84.74))

	// Run the actual endpoint
	router := mux.NewRouter()
	router.HandleFunc("/orders/{country}", c.CountOrders)
	router.ServeHTTP(w, r)
}
