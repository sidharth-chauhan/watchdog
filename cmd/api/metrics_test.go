package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMetricsEndpoint(t *testing.T) {
	// Create a new instance of our application
	app := newTestApplication(t)

	// Register the metric without starting the collection routine
	obaApiStatus.WithLabelValues(
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

	app := newTestApplication(t)

	testServer := ObaServer{
		Name:       "Test Server",
		ID:         999,
		ObaBaseURL: ts.URL,
		ObaApiKey:  "test-key",
	}

	// Test the checkServer function
	app.checkServer(testServer)

	// Wait a brief moment for metrics to be updated
	time.Sleep(100 * time.Millisecond)

	// Get and log all labels that are currently set for this metric
	metricChan := make(chan float64)
	go func() {
		metric, err := getMetricValue(obaApiStatus, map[string]string{
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
