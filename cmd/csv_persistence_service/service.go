package csvpersistenceservice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	controller "github.com/VictorPrado99/poc-csv-persistence/internal/csv_persistence_controller"
)

func StartService(port string) {
	router := controller.SetupRoutes()
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf(fmt.Sprintf("Starting Server on port %s \n", port))
	log.Fatal(srv.ListenAndServe())
}
