package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/OneBusAway/go-sdk/option"
	"github.com/julienschmidt/httprouter"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) pingMetricHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	server, exists := app.config.servers[id]

	if !exists {
		http.NotFound(w, r)
		return
	}

	client := onebusaway.NewClient(
		option.WithAPIKey(server.ObaApiKey),
		option.WithBaseURL(server.ObaBaseURL),
	)

	ctx := context.Background()

	currentTime, err := client.CurrentTime.Get(ctx)

	if err != nil {
		log.Fatalf("Error fetching current time: %v", err)
	}

	fmt.Fprintln(w, currentTime.Data.JSON.RawJSON())
}
