package main

import (
	service "github.com/VictorPrado99/poc-csv-persistence/cmd/csv_persistence_service"
	"github.com/VictorPrado99/poc-csv-persistence/internal/config"
	"github.com/VictorPrado99/poc-csv-persistence/internal/database"
)

func main() {
	// Load Configurations from config.json using Viper
	config.LoadAppConfig()
	// Initialize Database
	database.Connect(config.AppConfig.ConnectionString)
	database.Migrate()

	// Start api at specific port
	service.StartService(config.AppConfig.Port)
}
