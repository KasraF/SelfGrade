package main

import(
	"os"
	"net/http"
)

import (
	"fmt"
	"GoLog"
	"SelfGrade/code/persistance"
	"SelfGrade/code/web"
	"SelfGrade/code/security"
)

func main() {

	// Configure logger
	fmt.Println("Initializing logger...")
	GoLog.Init(os.Stdout)
	logger := GoLog.GetLogger()
	logger.Log("logger initialized.")

	// Configure Database
	logger.Log("Opening Database...")
	persistance.InitPostgreSQL()
	logger.Log("Database open.")

	// Configure security
	logger.Log("Initializing Security...")
	security.Init()
	logger.Log("Security initialized.")
	
	// Configure web handlers
	logger.Log("Configuring handlers...")
	web.Init()
	logger.Log("Handlers configured.")

	// Start server app
	logger.Log("Listening on Port 8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		logger.Error("Failed to listenAndServe():", err)
	}

	os.Exit(0)
}
