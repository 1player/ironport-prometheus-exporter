package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

var (
	listenAddress      = flag.String("web.listen-address", ":9101", "Listen address")
	metricsPath        = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	ironportHost       = flag.String("ironport.host", "localhost", "Hostname for the IronPort instance to monitor")
	ironportStatusPath = flag.String("ironport.status-path", "/xml/status", "Path to the IronPort XML status page")
)

func main() {
	flag.Parse()

	exporter, err := NewExporter(*ironportHost, *ironportStatusPath)
	if err != nil {
		log.Fatalln(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
