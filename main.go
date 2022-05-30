package main

import (
	service "github.com/VictorPrado99/poc-csv-persistence/cmd/csv_persistence_service"
	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
)

func main() {
	// Load Configurations from config.json using Viper
	LoadAppConfig()
	// Initialize Database
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()

	service.StartService(AppConfig.Port)
}
