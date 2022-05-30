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
)

// Method who will setup the router of this controller
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/orders", GetOrders).Methods(GET)
	r.HandleFunc("/orders/{id}", GetOrderById).Methods(GET)
	r.HandleFunc("/orders", CreateOrders).Methods(POST)
	r.HandleFunc("/orders/{id}", UpdateOrder).Methods(PUT)
	r.HandleFunc("/orders/{id}", DeleteOrder).Methods(DELETE)

	return r
}

func CreateOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/Orders", "POST Verb", "Called")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	var orders api.Orders

	json.NewDecoder(r.Body).Decode(&orders)
	for _, order := range orders {
		fmt.Println("Posting -> ", orders)
		database.Instance.Create(&order)
	}
	json.NewEncoder(w).Encode(orders)
}

func checkIfOrderExists(orderId string) bool {
	var order api.Order
	database.Instance.First(&order, orderId)

	return order.Id != 0
}

func GetOrderById(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	if !checkIfOrderExists(orderId) {
		json.NewEncoder(w).Encode("order Not Found!")
		return
	}
	var order api.Order
	database.Instance.First(&order, orderId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []api.Order

	fmt.Println("/Orders", "GET Verb", "Called")

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	fmt.Println("offset = ", offset)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	fmt.Println("limit = ", limit)

	sortParameter := r.URL.Query().Get("sort")
	if sortParameter == "" {
		sortParameter = "asc"
	} else if sortParameter != "asc" && sortParameter != "desc" {
		sortParameter = "asc"
	}

	fmt.Println("sortParameter = ", sortParameter)

	database.Instance.Offset(offset).Limit(limit).Order("id " + sortParameter).Find(&orders)

	fmt.Println("Retrieved data -> ", orders)

	// database.Instance.Find(&orders)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	if !checkIfOrderExists(orderId) {
		json.NewEncoder(w).Encode("order Not Found!")
		return
	}
	var order api.Order
	database.Instance.First(&order, orderId)
	json.NewDecoder(r.Body).Decode(&order)
	database.Instance.Save(&order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]
	if !checkIfOrderExists(orderId) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("order Not Found!")
		return
	}
	var order api.Order
	database.Instance.Delete(&order, orderId)
	json.NewEncoder(w).Encode("order Deleted Successfully!")
}
