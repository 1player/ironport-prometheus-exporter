package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/url"
)

const namespace = "ironport"

type Exporter struct {
	statusURL string

	up      prometheus.Gauge
	metrics map[string]prometheus.Gauge
}

func buildMetrics(names []string) map[string]prometheus.Gauge {
	m := make(map[string]prometheus.Gauge)
	for _, name := range names {
		m[name] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      name,
		})
	}

	return m
}

func NewExporter(host, statusPath string) (*Exporter, error) {
	statusURL := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   statusPath,
	}

	return &Exporter{
		statusURL: statusURL.String(),

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last Ironport scrape successful.",
		}),
		metrics: buildMetrics([]string{
			"inj_msgs",
			"inj_recips",
			"gen_bounce_recips",
			"rejected_recips",
			"dropped_msgs",
			"soft_bounced_evts",
			"completed_recips",
			"hard_bounced_recips",
			"dns_hard_bounced_recips",
			"5xx_hard_bounced_recips",
			"filter_hard_bounced_recips",
			"expired_hard_bounced_recips",
			"other_hard_bounced_recips",
			"delivered_recips",
			"deleted_recips",
			"global_unsub_hits",
			"ram_utilization",
			"total_utilization",
			"cpu_utilization",
			"av_utilization",
			"case_utilization",
			"bm_utilization",
			"disk_utilization",
			"resource_conservation",
			"log_used",
			"log_available",
			"conn_in",
			"conn_out",
			"active_recips",
			"unattempted_recips",
			"attempted_recips",
			"msgs_in_work_queue",
			"dests_in_memory",
			"kbytes_used",
			"kbytes_free",
			"msgs_in_policy_virus_outbreak_quarantine",
			"kbytes_in_policy_virus_outbreak_quarantine",
			"reporting_utilization",
			"quarantine_utilization",
		}),
	}, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape()

	ch <- e.up

	// Collect metrics
	for _, metric := range e.metrics {
		ch <- metric
	}
}

func (e *Exporter) scrape() {
	status, err := fetchStatus(e.statusURL)
	if err != nil {
		log.Println("Error scraping IronPort status page: ", err)
		e.up.Set(0)
		return
	}
	e.up.Set(1)

	// Process metrics
	for _, counter := range status.Counters {
		metric, ok := e.metrics[counter.Name]
		if !ok {
			continue
		}
		metric.Set(float64(counter.Lifetime))
	}
	for _, gauge := range status.Gauges {
		metric, ok := e.metrics[gauge.Name]
		if !ok {
			continue
		}
		metric.Set(float64(gauge.Current))
	}
}
