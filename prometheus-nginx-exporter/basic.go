package exporter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// NginxStats nginx basic stats
type NginxStats struct {
	// Nginx active connections
	ConnectionsActive float64

	// Connections (Reading - Writing - Waiting)
	Connections []Connections
}

type Connections struct {
	// Type is one of the (Reading - Writing - Waiting)
	Type string

	// Total number of connections
	Total float64
}

// scanBasicStats scans and parses nginx basic stats
func ScanBasicStats(r io.Reader) ([]NginxStats, error) {
	s := bufio.NewScanner(r)

	var stats []NginxStats
	var conns []Connections
	var nginxStats NginxStats

	for s.Scan() {
		fileds := strings.Fields(string(s.Bytes()))

		if len(fileds) == 3 && fileds[0] == "Active" {
			c, err := strconv.ParseFloat(fileds[2], 64)
			if err != nil {
				return nil, fmt.Errorf("%w: strconv.ParseFloat failed", err)
			}
			nginxStats.ConnectionsActive = c
		}

		if fileds[0] == "Reading:" {
			// fake metrics
			readingConns := Connections{Type: "reading", Total: 73}
			writingConns := Connections{Type: "writing", Total: 13}
			waitingConns := Connections{Type: "waiting", Total: 103}

			conns = append(conns, readingConns, writingConns, waitingConns)
			nginxStats.Connections = conns
		}

		// fmt.Println(fileds)

	}

	stats = append(stats, nginxStats)

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("%w, failed to read metrics", err)
	}

	return stats, nil

}
