package main

import (
	"time"

	"watchdog.onebusaway.org/internal/metrics"
)

func (app *application) startMetricsCollection() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, server := range app.config.Servers {
					metrics.ServerPing(server)
				}

				cachePath := "cache/gtfs.zip"
				_, _, err := metrics.CheckBundleExpiration(cachePath, app.logger)
				if err != nil {
					app.logger.Error("Failed to check GTFS bundle expiration", "error", err)
				}
			}
		}
	}()
}
