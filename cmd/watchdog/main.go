package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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

	serverName := flag.String("name", "", "Name of the OBA server")
	serverID := flag.Int("id", 0, "ID of the OBA server")
	baseURL := flag.String("base-url", "", "Base URL of the OBA server API")
	apiKey := flag.String("api-key", "", "API key for the OBA server")
	gtfsURL := flag.String("gtfs-url", "", "URL of the GTFS bundle")
	tripUpdateURL := flag.String("trip-update-url", "", "URL for trip updates")
	vehiclePositionURL := flag.String("vehicle-position-url", "", "URL for vehicle positions")

	flag.Parse()

	if *serverName == "" || *serverID == 0 || *baseURL == "" || *apiKey == "" || *gtfsURL == "" || *tripUpdateURL == "" || *vehiclePositionURL == "" {
		fmt.Println("Error: All flags are required.")
		flag.Usage()
		os.Exit(1)
	}

	server := models.NewObaServer(
		*serverName,
		*serverID,
		*baseURL,
		*apiKey,
		*gtfsURL,
		*tripUpdateURL,
		*vehiclePositionURL,
	)

	cfg.Servers = []models.ObaServer{*server}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Error("woot!!!")

	hash := sha1.Sum([]byte(*gtfsURL))
	hashStr := hex.EncodeToString(hash[:])

	cacheDir := "cache"
	cachePath := filepath.Join(cacheDir, hashStr+".zip")

	app := &application{
		config: cfg,
		logger: logger,
	}

	app.startMetricsCollection()

	// Download the GTFS bundle on startup
	err := utils.DownloadGTFSBundle(*gtfsURL, cachePath)
	if err != nil {
		logger.Error("Failed to download GTFS bundle", "error", err)
	} else {
		logger.Info("Successfully downloaded GTFS bundle", "path", cachePath)
	}

	// Cron job to download the GTFS bundle every 24 hours
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			err := utils.DownloadGTFSBundle(*gtfsURL, cachePath)
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
