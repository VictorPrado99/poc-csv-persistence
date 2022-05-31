package csvpersistencecontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
	"github.com/VictorPrado99/poc-csv-persistence/pkg/api"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
)

// Method who will setup the router of this controller
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/orders", GetOrders).Methods(GET)
	r.HandleFunc("/orders", CreateOrders).Methods(POST)
	r.HandleFunc("/orders/{country}", CountOrders).Methods(HEAD)

	return r
}

func CreateOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/Orders", "POST Verb", "Called")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var orders api.Orders

	json.NewDecoder(r.Body).Decode(&orders)
	database.Instance.CreateInBatches(orders, 3000)
	fmt.Println("Post Finished")
}

type HeadHeader struct {
	OrderCount int64
	Total      float64
}

func CountOrders(w http.ResponseWriter, r *http.Request) {
	country := mux.Vars(r)["country"]

	var data HeadHeader

	database.Instance.Model(&api.Order{}).Select("count(*) as order_count, sum(parcel_weight) as total").Where("country = ?", country).First(&data)

	w.Header().Add("x-weight-sum", fmt.Sprintf("%f", data.Total))
	w.Header().Add("x-orders-count", fmt.Sprintf("%d", data.OrderCount))
	w.Header().Add("x-country", country)

	w.WriteHeader(http.StatusOK)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []api.Order

	fmt.Println("/Orders", "GET Verb", "Called")

	query := r.URL.Query()

	offset, err := strconv.Atoi(query.Get("offset"))
	query.Del("offset")

	// Default Value
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	query.Del("limit")

	// Default Value
	if err != nil {
		limit = 10
	}

	// Hard limit
	if limit >= 100 {
		limit = 100
	}

	sortParameter := query.Get("sort")
	query.Del("sort")

	if sortParameter == "" {
		sortParameter = "asc"
	} else if sortParameter != "asc" && sortParameter != "desc" {
		sortParameter = "asc"
	}

	filterBy := "id " + sortParameter

	var scopesBuilder []func(db *gorm.DB) *gorm.DB

	countries := query.Get("country")

	if countries != "" {
		scopesBuilder = append(scopesBuilder, CountriesFilter(strings.Split(countries, ",")))
	}

	date := query.Get("date")

	if date != "" {
		scopesBuilder = append(scopesBuilder, DateFilter(date))
	}

	weightLimitStr := query.Get("weightLimit")

	if weightLimitStr != "" {
		weightLimitF, err := strconv.ParseFloat(weightLimitStr, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		weightLimit := float32(weightLimitF)
		scopesBuilder = append(scopesBuilder, WeightLimitFilter(weightLimit))
	}

	database.Instance.Offset(offset).Limit(limit).Order(filterBy).Scopes(scopesBuilder...).Find(&orders)

	if !(len(orders) > 0) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	nextOffset := offset + len(orders)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("x-next-offset", fmt.Sprintf("%d", nextOffset))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func CountriesFilter(countries []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country in (?)", countries)
	}
}

func DateFilter(date string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("date = (?)", date)
	}
}

func WeightLimitFilter(weightLimit float32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("parcel_weight <= (?)", weightLimit)
	}
}
