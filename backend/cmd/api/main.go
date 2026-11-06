package main

import (
	"log"

	"devtrack/internal/config"
	"devtrack/internal/database"
	"devtrack/internal/router"
)

func main() {
	//load the env to cfg variable
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config errors:%v", err)
	}
	//connect to DB
	dbPool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Config errors: %v", err)
		defer dbPool.Close()
	}
	//ini all the routes
	r := router.Setup(cfg) //pass creds to initialize the config

	//start the server
	addr := ":" + cfg.Port
	log.Printf("Server starting on port %s in %s mode...", cfg.Port, cfg.Environment)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to run:%v", err)
	}

}
