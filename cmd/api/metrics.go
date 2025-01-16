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

				numOfStaticAgencies, err := metrics.CheckAgenciesWithCoverage(cachePath, app.logger, app.config.Servers[0])
				numOfRealtimeAgencies, err := metrics.GetAgenciesWithCoverage(app.config.Servers[0])

				matchValue := 0
				if numOfRealtimeAgencies == numOfStaticAgencies {
					matchValue = 1
				}

				// 1 == match, 0 == no match
				metrics.AgenciesMatch.WithLabelValues(app.config.Servers[0].ObaBaseURL).Set(float64(matchValue))

				if err != nil {
					app.logger.Error("Failed to check agencies with coverage", "error", err)
				}
			}
		}
	}()
}
