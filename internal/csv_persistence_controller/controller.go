package csvpersistencecontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
	"github.com/VictorPrado99/poc-csv-persistence/pkg/api"
	"github.com/gorilla/mux"
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
	// r.HandleFunc("/orders/{id}", GetOrderById).Methods(GET)
	// r.HandleFunc("/orders/{id}", UpdateOrder).Methods(PUT)
	// r.HandleFunc("/orders/{id}", DeleteOrder).Methods(DELETE)

	return r
}

func CreateOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/Orders", "POST Verb", "Called")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	var orders api.Orders

	json.NewDecoder(r.Body).Decode(&orders)
	database.Instance.CreateInBatches(orders, 3000)
	fmt.Println("Post Finished")
}

// func checkIfOrderExists(orderId string) bool {
// 	var order api.Order
// 	database.Instance.First(&order, orderId)

// 	return order.Id != 0
// }

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

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

	// Default Value
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	// Default Value
	if err != nil {
		limit = 10
	}

	// Hard limit
	if limit >= 100 {
		limit = 100
	}

	sortParameter := r.URL.Query().Get("sort")
	if sortParameter == "" {
		sortParameter = "asc"
	} else if sortParameter != "asc" && sortParameter != "desc" {
		sortParameter = "asc"
	}

	filterBy := r.URL.Query().Get("filterBy")
	if filterBy == "" {
		filterBy = "id"
	}

	database.Instance.Offset(offset).Limit(limit).Order(filterBy + " " + sortParameter).Find(&orders)

	nextOffset := offset + len(orders)

	// database.Instance.Find(&orders)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("x-next-offset", fmt.Sprintf("%d", nextOffset))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// func GetOrderById(w http.ResponseWriter, r *http.Request) {
// 	orderId := mux.Vars(r)["id"]
// 	if !checkIfOrderExists(orderId) {
// 		json.NewEncoder(w).Encode("order Not Found!")
// 		return
// 	}
// 	var order api.Order
// 	database.Instance.First(&order, orderId)
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(order)
// }

// func UpdateOrder(w http.ResponseWriter, r *http.Request) {
// 	orderId := mux.Vars(r)["id"]
// 	if !checkIfOrderExists(orderId) {
// 		json.NewEncoder(w).Encode("order Not Found!")
// 		return
// 	}
// 	var order api.Order
// 	database.Instance.First(&order, orderId)
// 	json.NewDecoder(r.Body).Decode(&order)
// 	database.Instance.Save(&order)
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(order)
// }

// func DeleteOrder(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	orderId := mux.Vars(r)["id"]
// 	if !checkIfOrderExists(orderId) {
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode("order Not Found!")
// 		return
// 	}
// 	var order api.Order
// 	database.Instance.Delete(&order, orderId)
// 	json.NewEncoder(w).Encode("order Deleted Successfully!")
// }
