package main

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"log/slog"
	"os"
	"testing"
	"watchdog.onebusaway.org/internal/server"

	"watchdog.onebusaway.org/internal/models"
)

// getMetricValue is a helper function that retrieves the current value of a specific metric
func getMetricValue(metric *prometheus.GaugeVec, labels map[string]string) (float64, error) {
	// Create a collector for our specific metric
	c := make(chan prometheus.Metric, 1)
	metric.With(labels).Collect(c)

	// Get the metric from the channel
	m := <-c

	// Create a DESC and value for our metric
	var metricValue float64
	pb := &dto.Metric{}
	if err := m.Write(pb); err != nil {
		return 0, err
	}

	if pb.Gauge != nil {
		metricValue = pb.Gauge.GetValue()
	}

	return metricValue, nil
}

func newTestApplication(t *testing.T) *application {
	t.Helper()

	obaServer := models.NewObaServer(
		"Test Server",
		1,
		"https://test.example.com",
		"test-key",
		"",
		"",
		"",
	)

	cfg := server.NewConfig(
		4000,
		"testing",
		[]models.ObaServer{*obaServer},
	)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &application{
		config: *cfg,
		logger: logger,
	}
}

func getMetricsForTesting(t *testing.T, metric *prometheus.GaugeVec) {
	t.Helper()

	ch := make(chan prometheus.Metric)
	go func() {
		metric.Collect(ch)
		close(ch)
	}()

	for m := range ch {
		t.Logf("Found metric: %v", m.Desc())
	}
}
