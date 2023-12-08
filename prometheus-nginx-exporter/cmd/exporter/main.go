package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	exporter "github.com/mf210/custom-nginx-exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		targetHost = flag.String("target.host", "localhost", "nginx address with basic_status page")
		targetPort = flag.Int("target.port", 8080, "nginx port with basic_status page")
		targetPath = flag.String("target.path", "/status", "URL path to scrap metrics")
		promPort   = flag.Int("prom.port", 9150, "port to expose prometheus metrics")
	)
	flag.Parse()

	uri := fmt.Sprintf("http://%s:%d%s", *targetHost, *targetPort, *targetPath)

	// called on each collector.Collect.
	basicStatus := func() ([]exporter.NginxStats, error) {
		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := netClient.Get(uri)
		if err != nil {
			log.Fatalf("netClient.Get failed: %s: %s", uri, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("io.ReadAll failed: %s", err)
		}
		r := bytes.NewReader(body)

		return exporter.ScanBasicStats(r)
	}
	bc := exporter.NewBasicCollector(basicStatus)

	promRegis := prometheus.NewRegistry()
	promRegis.MustRegister(bc)
	// Go collector for test
	// promRegis.MustRegister(collectors.NewGoCollector())

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
