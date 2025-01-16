package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"watchdog.onebusaway.org/internal/metrics"
	"watchdog.onebusaway.org/internal/models"
)

func TestMetricsEndpoint(t *testing.T) {
	// Create a new instance of our application
	app := newTestApplication(t)

	// Register the metric without starting the collection routine
	metrics.ObaApiStatus.WithLabelValues(
		"1",
		"https://test.example.com",
	).Set(1)

	// Create a test server
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	// Make a request to the metrics endpoint
	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
	}

	// Check that the response contains our metric
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(body), "oba_api_status") {
		t.Error("metrics response doesn't contain oba_api_status metric")
	}
}

func TestCheckServer(t *testing.T) {
	requestChan := make(chan *http.Request, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Store the incoming request
		requestChan <- r

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":200,"currentTime":1234567890000,"text":"OK","version":2,"data":{"entry":{"readableTime":"Test Time"}}}`))
	}))
	defer ts.Close()

	testServer := models.ObaServer{
		Name:       "Test Server",
		ID:         999,
		ObaBaseURL: ts.URL,
		ObaApiKey:  "test-key",
	}

	// Test the checkServer function
	metrics.ServerPing(testServer)

	// Wait a brief moment for metrics to be updated
	time.Sleep(100 * time.Millisecond)

	// Get and log all labels that are currently set for this metric
	metricChan := make(chan float64)
	go func() {
		metric, err := getMetricValue(metrics.ObaApiStatus, map[string]string{
			"server_id":  "999",
			"server_url": testServer.ObaBaseURL,
		})
		if err != nil {
			t.Errorf("Failed to get metric value: %v", err)
		}
		//t.Logf("Got metric value: %v with labels server_id=999, server_url=%s",
		//	metric, testServer.ObaBaseURL)
		metricChan <- metric
	}()

	select {
	case metricValue := <-metricChan:
		if metricValue != 1 {
			t.Errorf("Expected metric value to be 1 (working), got %v", metricValue)
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for metric value")
	}
}

func TestCheckBundleExpiration(t *testing.T) {
	fixturePath, err := filepath.Abs(filepath.Join("..", "..", "testdata", "gtfs.zip"))
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	fixedTime := time.Date(2025, 1, 12, 20, 16, 38, 0, time.UTC)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	earliest, latest, err := metrics.CheckBundleExpiration(fixturePath, logger, fixedTime)
	if err != nil {
		t.Fatalf("CheckBundleExpiration failed: %v", err)
	}

	// This is the current gtfs.zip file in the testdata directory expected earliest and latest expiration days
	expectedEarliest := int(time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC).Sub(fixedTime).Hours() / 24)
	expectedLatest := int(time.Date(2025, 3, 28, 0, 0, 0, 0, time.UTC).Sub(fixedTime).Hours() / 24)

	if earliest != expectedEarliest {
		t.Errorf("Expected earliest expiration days to be %d, got %d", expectedEarliest, earliest)
	}
	if latest != expectedLatest {
		t.Errorf("Expected latest expiration days to be %d, got %d", expectedLatest, latest)
	}

	earliestMetric, err := getMetricValue(metrics.BundleEarliestExpirationGauge, map[string]string{
		"agency_id": "BundleExpiration",
	})
	if err != nil {
		t.Errorf("Failed to get earliest expiration metric value: %v", err)
	}
	if earliestMetric != float64(expectedEarliest) {
		t.Errorf("Expected earliest expiration metric to be %v, got %v", expectedEarliest, earliestMetric)
	}

	latestMetric, err := getMetricValue(metrics.BundleLatestExpirationGauge, map[string]string{
		"agency_id": "BundleExpiration",
	})
	if err != nil {
		t.Errorf("Failed to get latest expiration metric value: %v", err)
	}
	if latestMetric != float64(expectedLatest) {
		t.Errorf("Expected latest expiration metric to be %v, got %v", expectedLatest, latestMetric)
	}
}

func TestCheckAgenciesWithCoverage(t *testing.T) {
	fixturePath, err := filepath.Abs(filepath.Join("..", "..", "testdata", "gtfs.zip"))
	if err != nil {
		t.Fatalf("Failed to get absolute path to testdata/gtfs.zip: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	requestChan := make(chan *http.Request, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestChan <- r

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
            "code": 200,
            "currentTime": 1234567890000,
            "text": "OK",
            "version": 2,
            "data": {
                "list": [
                    {
                        "agencyId": "1"
                    }
                ]
            }
        }`))
	}))
	defer ts.Close()

	testServer := models.ObaServer{
		Name:       "Test Server",
		ID:         999,
		ObaBaseURL: ts.URL,
		ObaApiKey:  "test-key",
	}

	numOfStaticAgencies, err := metrics.CheckAgenciesWithCoverage(fixturePath, logger, testServer)
	if err != nil {
		t.Fatalf("CheckAgenciesWithCoverage failed: %v", err)
	}

	numOfRealtimeAgencies, err := metrics.GetAgenciesWithCoverage(testServer)
	if err != nil {
		t.Fatalf("GetAgenciesWithCoverage failed: %v", err)
	}

	matchValue := 0
	if numOfRealtimeAgencies == numOfStaticAgencies {
		matchValue = 1
	}
	metrics.AgenciesMatch.WithLabelValues(testServer.ObaBaseURL).Set(float64(matchValue))

	agencyMatchMetric, err := getMetricValue(metrics.AgenciesMatch, map[string]string{
		"server_id": testServer.ObaBaseURL,
	})
	if err != nil {
		t.Errorf("Failed to get AgenciesMatch metric value: %v", err)
	}

	if agencyMatchMetric != 1 {
		t.Errorf("Expected agency match metric to be 1, got %v", agencyMatchMetric)
	}

}
