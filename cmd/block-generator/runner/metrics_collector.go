package runner

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// MetricsCollector queries a /metrics endpoint for prometheus style metrics and saves metrics matching a pattern.
type MetricsCollector struct {
	// MetricsURL where metrics can be queried.
	MetricsURL string
	// Data is all of the results.
	Data []Entry
}

// Entry is the raw data pulled from the endpoint along with a timestamp.
type Entry struct {
	Timestamp time.Time
	Data      []string
}

// Collect fetches the metrics.
func (r *MetricsCollector) Collect(prefix string) error {
	metrics, err := r.getMetrics(prefix)
	if err != nil {
		return err
	}

	if len(metrics) > 0 {
		entry := Entry{Timestamp: time.Now()}
		for _, str := range metrics {
			entry.Data = append(entry.Data, str)
		}
		r.Data = append(r.Data, entry)
	}

	return nil
}

func (r MetricsCollector) getMetrics(prefix string) (result []string, err error) {
	resp, err := http.Get(r.MetricsURL)
	if err != nil {
		err = fmt.Errorf("unable to read metrics url '%s'", r.MetricsURL)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		str := scanner.Text()

		if strings.HasPrefix(str, "#") {
			continue
		}

		if strings.HasPrefix(str, prefix) {
			result = append(result, str)
		}
	}

	if scanner.Err() != nil {
		err = fmt.Errorf("problem reading metrics response: %w", scanner.Err())
	}

	return
}
