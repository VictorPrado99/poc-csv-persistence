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
	fmt.Println("[GET] /Orders", "Created")
	r.HandleFunc("/orders", CreateOrders).Methods(POST)
	fmt.Println("[POST] /Orders", "Created")
	r.HandleFunc("/orders/{country}", CountOrders).Methods(HEAD)
	fmt.Println("[HEAD] /Orders/{country}", "Created")

	return r
}

func CreateOrders(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "plain/text")
	var orders api.Orders

	json.NewDecoder(r.Body).Decode(&orders)
	err := database.Instance.CreateInBatches(orders, 3000).Error

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Println("Data Not Saved ->", err)
		fmt.Fprintf(w, "Couldn't save data ->  %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Println("Data Saved")
	fmt.Fprintf(w, "Data saved")
}

// Data type to use a single query to count
type HeadHeader struct {
	OrderCount int64
	Total      float64
}

func CountOrders(w http.ResponseWriter, r *http.Request) {
	country := mux.Vars(r)["country"] // Get the country passed as parameter

	var data HeadHeader

	// Query counting itens and making a sum of weight
	database.Instance.Model(&api.Order{}).Select("count(*) as order_count, sum(parcel_weight) as total").Where("country = ?", country).First(&data)

	// Prepare the header who will be the response
	w.Header().Add("x-weight-sum", fmt.Sprintf("%f", data.Total))
	w.Header().Add("x-orders-count", fmt.Sprintf("%d", data.OrderCount))
	w.Header().Add("x-country", country)

	// Send the header without a body
	w.WriteHeader(http.StatusOK)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []api.Order

	query := r.URL.Query() //Get the query object for simplicity sake

	// Offset data for pagination
	offset, err := strconv.Atoi(query.Get("offset"))

	// Default Value
	if err != nil { // If haven't receibed, consider it 0
		offset = 0
	}

	// Limit for page
	limit, err := strconv.Atoi(query.Get("limit"))

	// Default Value
	if err != nil { // If not received, consider it 10
		limit = 10
	}

	// Hard limit
	if limit >= 100 { // If the limit is too large, consider it 100
		limit = 100
	}

	sortParameter := query.Get("sort") // Get sort strategy

	if sortParameter == "" { // If no set, consider asc, and fill the variable just for consistency of query
		sortParameter = "asc"
	} else if sortParameter != "asc" && sortParameter != "desc" { //If something different from asc and desc, default to asc
		sortParameter = "asc"
	}

	filterBy := "id " + sortParameter

	// A slice to get function to create the WHERE Clause in Gorm, that way, we can build the clause dynamically base on functions
	var scopesBuilder []func(db *gorm.DB) *gorm.DB

	countries := query.Get("country") // Get contry or countries

	if countries != "" { // If have any country, split on commas, and generate a where using the IN clause
		scopesBuilder = append(scopesBuilder, CountriesFilter(strings.Split(countries, ","))) // Append the function in the slice
	}

	date := query.Get("date") //Get date

	if date != "" {
		scopesBuilder = append(scopesBuilder, DateFilter(date)) //append the function which the where clause for date is built
	}

	weightLimitStr := query.Get("weightLimit") // Get weight limit

	if weightLimitStr != "" {
		weightLimitF, err := strconv.ParseFloat(weightLimitStr, 32) // Parse to float64
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //If couldn't, is a bad request
			return
		}
		weightLimit := float32(weightLimitF)                                  //parse the float64 to float32
		scopesBuilder = append(scopesBuilder, WeightLimitFilter(weightLimit)) // append the function which the where clause is built
	}

	// Do the query, appending the offset, limit, spreading the functions which the where clause is built, and finally passing the wrapper struct, which is a slice of order struct
	database.Instance.Offset(offset).Limit(limit).Order(filterBy).Scopes(scopesBuilder...).Find(&orders)

	// If couldn't find anywhing, return Not Found
	if !(len(orders) > 0) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	nextOffset := offset + len(orders) // Set the offset of the next request to pass with the header.

	// set up the header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("x-next-offset", fmt.Sprintf("%d", nextOffset))
	w.WriteHeader(http.StatusOK)

	// Return the orders
	json.NewEncoder(w).Encode(orders)
}

//This function takes a slice of countries, and return a function thats return a gorm instance which contain a where clause for filter by this countries
func CountriesFilter(countries []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country in (?)", countries)
	}
}

// Filter by a specific date, and return a fuction which returns a gorm instance with this where clause
func DateFilter(date string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("date = (?)", date)
	}
}

// same as the two above, taking the weight limit and getting just what is equal or lesser
func WeightLimitFilter(weightLimit float32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("parcel_weight <= (?)", weightLimit)
	}
}
