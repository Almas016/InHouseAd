package cmd

import (
	"InHouseAd/internal/api"
	"InHouseAd/internal/api/handlers"
	"InHouseAd/internal/app"
	"log"
)

func Start() {
	wc := app.NewWebsiteChecker()

	// Load websites from url.txt file
	err := wc.LoadWebsitesFromFile("url.txt")
	if err != nil {
		log.Fatal("Failed to load websites:", err)
	}

	// Start the website availability checker
	go wc.CheckerWithTicker()

	// Create the web handler
	handler := handlers.NewHandler(wc)

	// Start the web server
	err = api.Routes(handler)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
