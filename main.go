package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"net/url"
)

var (
	listenAddress = flag.String("web.listen-address", ":9101", "Listen address")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")

	ironportHost           = flag.String("ironport.host", "localhost", "Hostname for the IronPort instance to monitor")
	ironportStatusPath     = flag.String("ironport.status-path", "/xml/status", "Path to the IronPort XML status page")
	ironportSkipCertVerify = flag.Bool("ironport.skip-cert-verify", false, "Skip SSL certificate verification")
	ironportAuth           = flag.String("ironport.basic-auth", "", "Ironport HTTP credentials in format user:pass")
)

func main() {
	flag.Parse()

	// Build status URL
	statusURL := url.URL{
		Scheme: "http",
		Host:   *ironportHost,
		Path:   *ironportStatusPath,
	}

	exporter, err := NewExporter(statusURL.String(), *ironportSkipCertVerify, *ironportAuth)
	if err != nil {
		log.Fatalln(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
