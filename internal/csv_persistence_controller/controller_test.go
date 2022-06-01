package csvpersistencecontroller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
	"github.com/VictorPrado99/poc-csv-persistence/pkg/api"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSetupRouter(t *testing.T) {
	router := SetupRoutes()

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

	CreateOrders(w, r)

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

	CreateOrders(w, r)

	if w.Result().StatusCode != 208 {
		t.Fatalf("Didn't receive status 208, received status %d", w.Result().StatusCode)
	}

	w.Result().Body.Close()
}

func TestHeadCount(t *testing.T) {
	// _, mock := StartMockMySqlDb()

	// r := httptest.NewRequest(http.MethodGet, "/orders", nil)

	// w := httptest.NewRecorder()
}

func TestGetWithoutFilters(t *testing.T) {
	_, mock := StartMockMySqlDb()

	orderType := reflect.TypeOf(api.Order{})
	var columns []string

	for i := 0; i < orderType.NumField(); i++ {
		columns = append(columns, orderType.Field(i).Name)
	}

	r := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` ORDER BY id asc LIMIT 10")).WillReturnRows(mock.NewRows(columns))

	GetOrders(w, r)

	data, _ := ioutil.ReadAll(w.Body)

	fmt.Println(string(data))
}
