package metrics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/OneBusAway/go-sdk/option"
	"github.com/jamespfennell/gtfs"
	"watchdog.onebusaway.org/internal/models"
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

	RealtimeVehiclePositions.WithLabelValues(gtfsRtURL).Set(float64(count))

	return count, nil
}

func VehiclesForAgencyAPI(server models.ObaServer) (int, error) {

	client := onebusaway.NewClient(
		option.WithAPIKey(server.ObaApiKey),
		option.WithBaseURL(server.ObaBaseURL),
	)

	ctx := context.Background()

	agencyID := "unitrans"

	response, err := client.VehiclesForAgency.List(ctx, agencyID, onebusaway.VehiclesForAgencyListParams{})

	if err != nil {
		return 0, err
	}

	VehicleCountAPI.WithLabelValues(agencyID).Set(float64(len(response.Data.List)))

	return len(response.Data.List), nil
}

func CheckVehicleCountMatch(vehiclePositionsURL, apiKey, apiValue string, server models.ObaServer) error {
	gtfsRtVehicleCount, err := CountVehiclePositions(vehiclePositionsURL, apiKey, apiValue)
	if err != nil {
		return fmt.Errorf("failed to count vehicle positions from GTFS-RT: %v", err)
	}

	apiVehicleCount, err := VehiclesForAgencyAPI(server)
	if err != nil {
		return fmt.Errorf("failed to count vehicle positions from API: %v", err)
	}

	match := 0
	if gtfsRtVehicleCount == apiVehicleCount {
		match = 1
	}

	VehicleCountMatch.WithLabelValues(server.ObaBaseURL).Set(float64(match))

	return nil
}
