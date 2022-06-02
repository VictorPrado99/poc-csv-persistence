package csvpersistencecontroller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/VictorPrado99/poc-csv-persistence/internal/config"
	c "github.com/VictorPrado99/poc-csv-persistence/internal/csv_persistence_controller"
	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
	"github.com/VictorPrado99/poc-csv-persistence/pkg/api"
)

const (
	COMPOSE_PATH = "../../dev/docker-compose.yml"
)

var (
	// Generate Test Data
	orderTestData = api.Orders{
		{
			Id:           1,
			Email:        "test1@test.com",
			PhoneNumber:  "256 11 6583 2753",
			ParcelWeight: 28.46,
			Date:         "2022-05-31",
			Country:      "Uganda",
		},
		{
			Id:           2,
			Email:        "test2@test.com",
			PhoneNumber:  "055 11 6583 2753",
			ParcelWeight: 54.74,
			Date:         "2022-05-31",
			Country:      "Brazil",
		},
		{
			Id:           3,
			Email:        "test3@test.com",
			PhoneNumber:  "351 11 6583 2753",
			ParcelWeight: 75.14,
			Date:         "2022-06-20",
			Country:      "Portugal",
		},
		{
			Id:           4,
			Email:        "test4@test.com",
			PhoneNumber:  "055 11 6583 2753",
			ParcelWeight: 30,
			Date:         "2022-06-20",
			Country:      "Brazil",
		},
		{
			Id:           5,
			Email:        "test5@test.com",
			PhoneNumber:  "351 11 6583 2753",
			ParcelWeight: 70,
			Date:         "2022-05-31",
			Country:      "Portugal",
		},
		{
			Id:           6,
			Email:        "test6@test.com",
			PhoneNumber:  "351 11 6583 2753",
			ParcelWeight: 29.13,
			Date:         "2022-05-31",
			Country:      "Portugal",
		},
	}
)

// The test Suite fixture
type GetOrdersTestSuite struct {
	suite.Suite
}

// This is not the best solution for test. I made using the compose, to faster development, but the right would be each test suite running his own infrasctructe without compose
// Two tests suites using databases woudn't be possible that way, because of port listening
func (suite *GetOrdersTestSuite) ComposeUp() error {
	time.Sleep(2 * time.Second) //For safety, sometimes between tests, they tried to start before properly ended

	composeFilePaths := []string{COMPOSE_PATH}
	identifier := strings.ToLower(uuid.New().String())

	compose := tc.NewLocalDockerCompose(composeFilePaths, identifier)

	// Set compose down after ending the suit
	suite.T().Cleanup(func() {
		compose.Down()
	})

	execError := compose.
		WithCommand([]string{"up", "-d"}).
		WaitForService(identifier+"_db_1", wait.ForLog("port: 3306  MySQL Community Server - GPL")).
		Invoke()
	err := execError.Error

	if err != nil {
		return fmt.Errorf("Could not run compose file: %v - %v", composeFilePaths, err)
	}

	return nil
}

// Setup the test suite fixture
func (suite *GetOrdersTestSuite) SetupTest() {
	config.LoadAppConfig() //Get config file

	// Compose up, and set to be down at cleanup
	err := suite.ComposeUp()

	if err != nil {
		suite.T().Fatalf("Problem setting infrastructure")
	}

	// Connect to created database
	database.Connect(config.AppConfig.ConnectionString)
	database.Migrate()

	database.Instance.CreateInBatches(&orderTestData, 50)

}

// Run tests in the suite
func TestGetOrdersSuite(t *testing.T) {
	testTable = map[string]func() api.Orders{
		"/orders":                       func() api.Orders { return orderTestData },
		"/orders?offset=4":              generateOffset4Data,
		"/orders?limit=3":               generateLimit3Data,
		"/orders?sort=desc":             generateSortedDescData,
		"/orders?country=Uganda":        generateFilterByCountryUganda,
		"/orders?country=Uganda,Brazil": generateFilterByCountryUgandaBrazil,
		"/orders?date=2022-06-20":       generateFilterByDate0620,
		"/orders?weightLimit=30.1":      generateFilter30LessEqual,
	}

	suite.Run(t, new(GetOrdersTestSuite))
}

var testTable map[string]func() api.Orders

func (suite *GetOrdersTestSuite) TestGetOrders() {

	for path, getReturn := range testTable {
		func(path string, getReturn func() api.Orders) {
			suite.T().Logf("Running test for endpoint -> %s", path)

			// Generate mock for request and write
			r := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()

			c.GetOrders(w, r)

			var got api.Orders

			err := json.Unmarshal(w.Body.Bytes(), &got)

			if err != nil {
				suite.T().Errorf("Error %v, when converting the body", err)
			}

			want := getReturn()

			if !assert.Equal(suite.T(), want, got) {
				suite.T().Errorf("\nFailed when comparing results for endpoint %s \n", path)
			} else {
				suite.T().Logf("Success when comparing results for endpoint %s", path)
			}

		}(path, getReturn)
	}
}

func generateSortedDescData() api.Orders {
	var descData api.Orders
	lastIndex := len(orderTestData) - 1
	for i, _ := range orderTestData {
		descData = append(descData, orderTestData[lastIndex-i])
	}
	return descData
}

func generateOffset4Data() api.Orders {
	return orderTestData[4:]
}

func generateLimit3Data() api.Orders {
	return orderTestData[:3]
}

func generateFilterByCountryUganda() api.Orders {
	return api.Orders{
		{
			Id:           1,
			Email:        "test1@test.com",
			PhoneNumber:  "256 11 6583 2753",
			ParcelWeight: 28.46,
			Date:         "2022-05-31",
			Country:      "Uganda",
		},
	}
}

func generateFilterByCountryUgandaBrazil() api.Orders {
	return api.Orders{
		{
			Id:           1,
			Email:        "test1@test.com",
			PhoneNumber:  "256 11 6583 2753",
			ParcelWeight: 28.46,
			Date:         "2022-05-31",
			Country:      "Uganda",
		},
		{
			Id:           2,
			Email:        "test2@test.com",
			PhoneNumber:  "055 11 6583 2753",
			ParcelWeight: 54.74,
			Date:         "2022-05-31",
			Country:      "Brazil",
		},
		{
			Id:           4,
			Email:        "test4@test.com",
			PhoneNumber:  "055 11 6583 2753",
			ParcelWeight: 30,
			Date:         "2022-06-20",
			Country:      "Brazil",
		},
	}
}

func generateFilterByDate0620() api.Orders {
	return api.Orders{
		{
			Id:           3,
			Email:        "test3@test.com",
			PhoneNumber:  "351 11 6583 2753",
			ParcelWeight: 75.14,
			Date:         "2022-06-20",
			Country:      "Portugal",
		},
		{
			Id:           4,
			Email:        "test4@test.com",
			PhoneNumber:  "055 11 6583 2753",
			ParcelWeight: 30,
			Date:         "2022-06-20",
			Country:      "Brazil",
		},
	}
}

func generateFilter30LessEqual() api.Orders {
	return api.Orders{
		{
			Id:           1,
			Email:        "test1@test.com",
			PhoneNumber:  "256 11 6583 2753",
			ParcelWeight: 28.46,
			Date:         "2022-05-31",
			Country:      "Uganda",
		},
		{
			Id:           4,
			Email:        "test4@test.com",
			PhoneNumber:  "055 11 6583 2753",
			ParcelWeight: 30,
			Date:         "2022-06-20",
			Country:      "Brazil",
		},
		{
			Id:           6,
			Email:        "test6@test.com",
			PhoneNumber:  "351 11 6583 2753",
			ParcelWeight: 29.13,
			Date:         "2022-05-31",
			Country:      "Portugal",
		},
	}
}

func (suite *GetOrdersTestSuite) TestHeadOrdersBrazil() {
	path := "/orders/Brazil"

	suite.T().Logf("requesting head at  -> %s", path)

	// Generate mock for request and write
	r := httptest.NewRequest(http.MethodHead, path, nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/orders/{country}", c.CountOrders)
	router.ServeHTTP(w, r)

	country := w.Result().Header.Get("x-country")
	ordersCount := w.Result().Header.Get("X-Orders-Count")
	weightSum := w.Result().Header.Get("X-Weight-Sum")

	if !(country == "Brazil" &&
		ordersCount == "2" &&
		weightSum == "84.740005") {
		suite.T().Errorf("Header different than expected\n X-Country: %s \n X-Orders-Count: %s \n X=Weight-Sum: %s", country, ordersCount, weightSum)
	}
}
