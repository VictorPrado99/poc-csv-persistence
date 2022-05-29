package csvpersistencecontroller

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"
)

// Method who will setup the router of this controller
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/orders", postOrder).Methods(POST)

	return r
}

func postOrder(http.ResponseWriter, *http.Request) {
	// TODO Call the service to make the CRUD
}
