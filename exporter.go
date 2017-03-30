package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "ironport"

type Exporter struct {
	host, statusPath string

	up prometheus.Gauge
}

func NewExporter(host, statusPath string) (*Exporter, error) {
	return &Exporter{
		host:       host,
		statusPath: statusPath,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "up",
			Help: "Was the last Ironport scrape successful.",
		}),
	}, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

}
