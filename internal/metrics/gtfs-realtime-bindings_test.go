package metrics

import (
	"net/http"
	"testing"

	"watchdog.onebusaway.org/internal/models"
)

func TestCountVehiclePositions(t *testing.T) {
	t.Run("Valid GTFS-RT response", func(t *testing.T) {
		mockServer := setupGtfsRtServer(t, "gtfs_rt_feed_vehicles.pb")
		defer mockServer.Close()

		server := models.ObaServer{
			ID:                 1,
			VehiclePositionUrl: mockServer.URL,
			GtfsRtApiKey:       "Authorization",
			GtfsRtApiValue:     "test-key",
		}

		count, err := CountVehiclePositions(server)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if count < 0 {
			t.Fatalf("Expected non-negative count, got %d", count)
		}
	})

	t.Run("Unreachable server", func(t *testing.T) {
		server := models.ObaServer{
			ID:                 3,
			VehiclePositionUrl: "http://nonexistent.local/gtfs-rt",
		}

		_, err := CountVehiclePositions(server)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		server := models.ObaServer{
			ID:                 4,
			VehiclePositionUrl: "://invalid-url",
		}

		_, err := CountVehiclePositions(server)
		if err == nil {
			t.Fatal("Expected an error due to invalid URL, got nil")
		}
	})
}

func TestVehiclesForAgencyAPI(t *testing.T) {
	t.Run("NilResponse", func(t *testing.T) {
			ts := setupObaServer(t, `{"data": {"list": []}}`, http.StatusOK)
			defer ts.Close()

			server := models.ObaServer{
					Name:       "Test Server",
					ID:         999,
					ObaBaseURL: ts.URL,
					ObaApiKey:  "test-key",
					AgencyID:   "test-agency",
			}

			count, err := VehiclesForAgencyAPI(server)
			if err != nil {
					t.Fatalf("Expected no error, got %v", err)
			}

			if count != 0 {
					t.Fatalf("Expected count to be 0, got %d", count)
			}
	})

	t.Run("SuccessfulResponse", func(t *testing.T) {
			ts := setupObaServer(t, `{"data": {"list": [{"vehicleId": "1"}, {"vehicleId": "2"}]}}`, http.StatusOK)
			defer ts.Close()

			server := models.ObaServer{
					Name:       "Test Server",
					ID:         999,
					ObaBaseURL: ts.URL,
					ObaApiKey:  "test-key",
					AgencyID:   "test-agency",
			}

			count, err := VehiclesForAgencyAPI(server)
			if err != nil {
					t.Fatalf("Expected no error, got %v", err)
			}

			if count != 2 {
					t.Fatalf("Expected count to be 2, got %d", count)
			}
	})

	t.Run("ErrorResponse", func(t *testing.T) {
			ts := setupObaServer(t, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
			defer ts.Close()

			server := models.ObaServer{
					Name:       "Test Server",
					ID:         999,
					ObaBaseURL: ts.URL,
					ObaApiKey:  "test-key",
					AgencyID:   "test-agency",
			}

			_, err := VehiclesForAgencyAPI(server)
			if err == nil {
					t.Fatal("Expected an error but got nil")
			}
	})
}