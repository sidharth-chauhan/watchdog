package main

import (
	"time"

	"github.com/joho/godotenv"
	"watchdog.onebusaway.org/internal/metrics"
	"watchdog.onebusaway.org/internal/utils"
)

func (app *application) startMetricsCollection() {

	err := godotenv.Load()

	if err != nil {
		app.logger.Error("Failed to load .env file", "error", err)
	}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:

				app.mu.Lock()
				servers := app.config.Servers
				app.mu.Unlock()

				for _, server := range servers {
					metrics.ServerPing(server)
				}

				cachePath, err := utils.GetLastCachedFile("cache")

				if err != nil {
					app.logger.Error("Failed to get last cached file", "error", err)
					continue
				}

				_, _, err = metrics.CheckBundleExpiration(cachePath, app.logger, time.Now())
				if err != nil {
					app.logger.Error("Failed to check GTFS bundle expiration", "error", err)
				}

				err = metrics.CheckAgenciesWithCoverageMatch(cachePath, app.logger, app.config.Servers[0])

				if err != nil {
					app.logger.Error("Failed to check agencies with coverage match metric", "error", err)
				}

				err = metrics.CheckVehicleCountMatch(app.config.Servers[0].VehiclePositionUrl, app.config.Servers[0].GtfsRtApiKey, app.config.Servers[0].GtfsRtApiValue, app.config.Servers[0])

				if err != nil {
					app.logger.Error("Failed to check vehicle count match metric", "error", err)
				}
			}
		}
	}()
}
