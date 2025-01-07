package main

import (
	"context"
	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/OneBusAway/go-sdk/option"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"

	"watchdog.onebusaway.org/internal/models"
)

var (
	// API Status (up/down)
	obaApiStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "oba_api_status",
			Help: "Status of the OneBusAway API Server (0 = not working, 1 = working)",
		},
		[]string{"server_id", "server_url"},
	)
)

func (app *application) startMetricsCollection() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, server := range app.config.Servers {
					app.checkServer(server)
				}
			}
		}
	}()
}

func (app *application) checkServer(server models.ObaServer) {
	client := onebusaway.NewClient(
		option.WithAPIKey(server.ObaApiKey),
		option.WithBaseURL(server.ObaBaseURL),
	)

	ctx := context.Background()
	response, err := client.CurrentTime.Get(ctx)

	if err != nil {
		// Update status metric
		obaApiStatus.WithLabelValues(
			strconv.Itoa(server.ID),
			server.ObaBaseURL,
		).Set(0)
		return
	}

	// Check response validity
	if response.Data.Entry.ReadableTime != "" {
		obaApiStatus.WithLabelValues(
			strconv.Itoa(server.ID),
			server.ObaBaseURL,
		).Set(1)
	} else {
		obaApiStatus.WithLabelValues(
			strconv.Itoa(server.ID),
			server.ObaBaseURL,
		).Set(0)
	}
}
