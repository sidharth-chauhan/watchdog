package main

import (
	"time"

	"watchdog.onebusaway.org/internal/metrics"
	"watchdog.onebusaway.org/internal/utils"
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

				cachePath, err := utils.GetLastCachedFile("cache")

				if err != nil {
					app.logger.Error("Failed to get last cached file", "error", err)
					continue
				}

				_, _, err = metrics.CheckBundleExpiration(cachePath, app.logger, time.Now())
				if err != nil {
					app.logger.Error("Failed to check GTFS bundle expiration", "error", err)
				}

				staticAgenciesNumber, err := metrics.CheckAgenciesWithCoverage(cachePath, app.logger, app.config.Servers[0])
				if err != nil {
					app.logger.Error("Failed to check agencies with coverage", "error", err)
					continue
				}

				numberOfAgencies, err := metrics.GetAgenciesWithCoverage(app.config.Servers[0])
				if err != nil {
					app.logger.Error("Failed to get agencies with coverage", "error", err)
					continue
				}

				difference := staticAgenciesNumber - numberOfAgencies

				metrics.AgenciesDifference.WithLabelValues(app.config.Servers[0].ObaBaseURL).Set(float64(difference))

				// to verify that the number of agencies in the static GTFS file matches the number of agencies in the agencies-with-coverage endpoint
				// use this query in Prometheus: oba_agencies_difference == bool 0
				metrics.AgenciesDifference.WithLabelValues(app.config.Servers[0].ObaBaseURL).Set(float64(difference))
			}
		}
	}()
}
