package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		promPort = flag.Int("prom.port", 9150, "port to expose prometheus metrics")
	)
	flag.Parse()

	promRegis := prometheus.NewRegistry()
	// Go collector for test
	promRegis.MustRegister(collectors.NewGoCollector())

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(promRegis, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	// start listening for HTTP connections.
	port := fmt.Sprintf(":%d", *promPort)
	log.Printf("starting nginx exporter at on %q/metrics", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("cannot start nginx exporter: %s", err)
	}

}
