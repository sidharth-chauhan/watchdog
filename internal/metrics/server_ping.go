package metrics

import (
	"context"
	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/OneBusAway/go-sdk/option"
	"strconv"
	"watchdog.onebusaway.org/internal/models"
)

func ServerPing(server models.ObaServer) {
	client := onebusaway.NewClient(
		option.WithAPIKey(server.ObaApiKey),
		option.WithBaseURL(server.ObaBaseURL),
	)

	ctx := context.Background()
	response, err := client.CurrentTime.Get(ctx)

	if err != nil {
		// Update status metric
		ObaApiStatus.WithLabelValues(
			strconv.Itoa(server.ID),
			server.ObaBaseURL,
		).Set(0)
		return
	}

	// Check response validity
	if response.Data.Entry.ReadableTime != "" {
		ObaApiStatus.WithLabelValues(
			strconv.Itoa(server.ID),
			server.ObaBaseURL,
		).Set(1)
	} else {
		ObaApiStatus.WithLabelValues(
			strconv.Itoa(server.ID),
			server.ObaBaseURL,
		).Set(0)
	}
}
