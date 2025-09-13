package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("=== SHELLY EXPORTER DEBUG BUILD ===")
	config := getConfig()
	log.Printf("Loaded config: %+v\n", config)
	registerMetrics()
	log.Println("Metrics registered.")

	go func(config configuration) {
		log.Println("Starting device polling goroutine...")
		fetchDevices(config)
		for range time.Tick(config.ScrapeInterval) {
			log.Println("Polling devices...")
			fetchDevices(config)
		}
	}(config)

	log.Printf("starting web server on port %d", config.Port)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
