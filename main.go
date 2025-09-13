package main

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
       debug := os.Getenv("DEBUG") != ""
       if debug {
	       log.Println("=== SHELLY EXPORTER DEBUG BUILD ===")
       }
       config := getConfig()
       if debug {
	       log.Printf("Loaded config: %+v\n", config)
       }
       registerMetrics()
       if debug {
	       log.Println("Metrics registered.")
       }

       go func(config configuration) {
	       if debug {
		       log.Println("Starting device polling goroutine...")
	       }
	       fetchDevices(config)
	       for range time.Tick(config.ScrapeInterval) {
		       if debug {
			       log.Println("Polling devices...")
		       }
		       fetchDevices(config)
	       }
       }(config)

       log.Printf("starting web server on port %d", config.Port)
       http.Handle("/metrics", promhttp.Handler())
       http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
