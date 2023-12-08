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
	ConnectionsActive float64
}

// scanBasicStats scans and parses nginx basic stats
func ScanBasicStats(r io.Reader) ([]NginxStats, error) {
	s := bufio.NewScanner(r)

	var stats []NginxStats
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
		// fmt.Println(fileds)

	}

	stats = append(stats, nginxStats)

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("%w, failed to read metrics", err)
	}

	return stats, nil

}
