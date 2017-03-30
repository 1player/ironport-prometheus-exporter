package main

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

// Counter is a counter returned in the XML status document
type Counter struct {
	Name     string `xml:"name,attr"`
	Lifetime uint64 `xml:"lifetime,attr"`
}

// Gauge is a gauge returned in the XML status document
type Gauge struct {
	Name    string     `xml:"name,attr"`
	Current GaugeValue `xml:"current,attr"`
}

// SystemStatus is the system status returned in the XML status document
type SystemStatus struct {
	Status string `xml:"status,attr"`
}

// Status is the root of the XML status document returned by Ironport
type Status struct {
	System   SystemStatus `xml:"system"`
	Counters []Counter    `xml:"counters>counter"`
	Gauges   []Gauge      `xml:"gauges>gauge"`
}

type GaugeValue int64

var gaugeValueRegex = regexp.MustCompile(`([0-9]+)([KMGT]?)`)

func (g *GaugeValue) UnmarshalXMLAttr(attr xml.Attr) error {
	found := gaugeValueRegex.FindStringSubmatch(attr.Value)
	if found == nil {
		return fmt.Errorf("Invalid gauge value %v", attr.Value)
	}

	// Parse base value
	v, err := strconv.ParseInt(found[1], 10, 64)
	if err != nil {
		return err
	}

	// Apply multiplier
	switch found[2] {
	case "K":
		v *= 1024
	case "M":
		v *= 1024 * 1024
	case "G":
		v *= 1024 * 1024 * 1024
	case "T":
		v *= 1024 * 1024 * 1024 * 1024
	}

	*g = GaugeValue(v)
	return nil
}

func parseStatus(r io.Reader) (*Status, error) {
	var status Status

	decoder := xml.NewDecoder(r)
	err := decoder.Decode(&status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func fetchStatus(url string, skipCertVerify bool, user, pass string) (*Status, error) {
	// Build client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *ironportSkipCertVerify},
	}
	client := &http.Client{Transport: tr}

	// Build request
	req, err := http.NewRequest("GET", url, http.NoBody)
	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	// Fetch response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseStatus(resp.Body)
}
