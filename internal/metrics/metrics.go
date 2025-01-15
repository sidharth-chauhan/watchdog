package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ObaApiStatus API Status (up/down)
	ObaApiStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "oba_api_status",
			Help: "Status of the OneBusAway API Server (0 = not working, 1 = working)",
		},
		[]string{"server_id", "server_url"},
	)
)

var (
	BundleEarliestExpirationGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gtfs_bundle_days_until_earliest_expiration",
		Help: "Number of days until the earliest GTFS bundle expiration",
	}, []string{"agency_id"})

	BundleLatestExpirationGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gtfs_bundle_days_until_latest_expiration",
		Help: "Number of days until the latest GTFS bundle expiration",
	}, []string{"agency_id"})
)

var (
	AgenciesInStaticGtfs = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "oba_agencies_in_static_gtfs",
		Help: "Number of agencies in the static GTFS file",
	}, []string{"server_id"})

	AgenciesInCoverageEndpoint = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "oba_agencies_in_coverage_endpoint",
		Help: "Number of agencies in the agencies-with-coverage endpoint",
	}, []string{"server_id"})

	AgenciesMatch = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "oba_agencies_match",
		Help: "Whether the number of agencies in the static GTFS file matches the agencies-with-coverage endpoint (1 = match, 0 = no match)",
	}, []string{"server_id"})
)
