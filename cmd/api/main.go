package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"watchdog.onebusaway.org/internal/models"
	"watchdog.onebusaway.org/internal/server"
	"watchdog.onebusaway.org/internal/utils"
)

// Declare a string containing the application version number. Later in the book we'll
// generate this automatically at build time, but for now we'll just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config server.Config
	logger *slog.Logger
}

func main() {
	var cfg server.Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	server := models.NewObaServer(
		"Sound Transit",
		1,
		"https://api.pugetsound.onebusaway.org",
		"org.onebusaway.iphone",
		"https://www.soundtransit.org/GTFS-rail/40_gtfs.zip",
		"https://api.pugetsound.onebusaway.org/api/gtfs_realtime/trip-updates-for-agency/40.pb?key=org.onebusaway.iphone",
		"https://api.pugetsound.onebusaway.org/api/gtfs_realtime/vehicle-positions-for-agency/40.pb?key=org.onebusaway.iphone",
	)

	cfg.Servers = []models.ObaServer{*server}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	app.startMetricsCollection()

	// Download the GTFS bundle on startup
	cachePath := "cache/gtfs.zip"
	bundleURL := "https://www.soundtransit.org/GTFS-rail/40_gtfs.zip"
	err := utils.DownloadGTFSBundle(bundleURL, cachePath)
	if err != nil {
		logger.Error("Failed to download GTFS bundle", "error", err)
	} else {
		logger.Info("Successfully downloaded GTFS bundle", "path", cachePath)
	}

	// Cron job to download the GTFS bundle every 24 hours
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			err := utils.DownloadGTFSBundle(bundleURL, cachePath)
			if err != nil {
				logger.Error("Failed to download GTFS bundle", "error", err)
			} else {
				logger.Info("Successfully updated GTFS bundle", "path", cachePath)
			}
		}
	}()

	// Use the httprouter instance returned by app.routes() as the server handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
