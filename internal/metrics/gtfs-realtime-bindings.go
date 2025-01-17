package metrics

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jamespfennell/gtfs"
)

func CountVehiclePositions(gtfsRtURL string, apiKey string, apiValue string) (int, error) {
	parsedURL, err := url.Parse(gtfsRtURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse GTFS-RT URL: %v", err)
	}

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set(apiKey, apiValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch GTFS-RT feed: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read GTFS-RT feed: %v", err)
	}

	realtimeData, err := gtfs.ParseRealtime(data, &gtfs.ParseRealtimeOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to parse GTFS-RT feed: %v", err)
	}

	count := len(realtimeData.Vehicles)

	RealtimeVehiclePositions.Set(float64(count))

	return count, nil
}
