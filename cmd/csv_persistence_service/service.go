package csvpersistenceservice

import (
	"log"
	"net/http"
	"time"

	controller "github.com/VictorPrado99/poc-csv-persistence/internal/csv_persistence_controller"
)

func StartService() {
	router := controller.SetupRoutes()
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9100",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
