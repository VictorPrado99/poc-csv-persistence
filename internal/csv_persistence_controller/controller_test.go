package csvpersistencecontroller

import (
	"bytes"
	"database/sql"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSetupRouter(t *testing.T) {
	router := SetupRoutes()

	if router == nil {
		t.Fatalf("Failed to create router")
	}
}

func StartMockMySqlDb(db *sql.DB, mock sqlmock.Sqlmock) {
	dialec := mysql.New(mysql.Config{
		Conn: db,
	})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).WillReturnRows(mock.NewRows([]string{"version"}).AddRow(1))
	database.Instance, _ = gorm.Open(dialec, &gorm.Config{})
}

func TestCreateOrders(t *testing.T) {
	db, mock, _ := sqlmock.New()

	StartMockMySqlDb(db, mock)

	str := `{"id": 1,"email": "test1e@teste","phone_number": "351 961 251 326","parcel_weight": 22.4,"date": "2022-03-12","country": "Portugal"}`

	dataMock := []byte(str)

	r := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(dataMock))

	w := httptest.NewRecorder()

	mock.ExpectBegin()
	mock.ExpectCommit()
	// mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `orders` (`email`,`phone_number`,`parcel_weight`,`date`,`country`,`id`) VALUES (?)")).
	// 	WithArgs("test1e@teste", "351 961 251 326", 22.4, "2022-03-12", "Portugal").
	// 	WillReturnResult(sqlmock.NewResult(0, 1))

	// mock.ExpectClose()

	CreateOrders(w, r)

	// defer w.Result().Body.Close()

	// data, err := ioutil.ReadAll(r.Body)

	// if err != nil {
	// 	t.Errorf("expected error to be nil got %v", err)
	// }

	// fmt.Println(string(data))

	// if string(data) ==

	// var data struct{}

	// json.Unmarshal(dataMock, &data)

	// fmt.Println(data)

	// mock.ExpectBegin()
	// mock.ExpectQuery("INSERT INTO `orders` (`email`,`phone_number`,`parcel_weight`,`date`,`country`,`id`) VALUES ('test1e@teste','351 961 251 326',26.400000,'2022-03-12','Portugal',100065421),('test1e@teste','351 961 251 326',0.000000,'2022-03-12','Portugal',100065422),('test1e@teste','351 961 251 326',0.000000,'2022-03-12','Portugal',100065421),('test1e@teste','351 961 251 326',0.000000,'2022-03-12','Portugal',4)")
}
