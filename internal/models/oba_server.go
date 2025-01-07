package models

// ObaServer represents a OneBusAway server configuration
type ObaServer struct {
	Name               string
	ID                 int
	ObaBaseURL         string
	ObaApiKey          string
	GtfsUrl            string
	TripUpdateUrl      string
	VehiclePositionUrl string
}

// NewObaServer creates a new ObaServer instance with the provided configuration
func NewObaServer(name string, id int, baseURL, apiKey, gtfsURL, tripUpdateURL, vehiclePositionURL string) *ObaServer {
	return &ObaServer{
		Name:               name,
		ID:                 id,
		ObaBaseURL:         baseURL,
		ObaApiKey:          apiKey,
		GtfsUrl:            gtfsURL,
		TripUpdateUrl:      tripUpdateURL,
		VehiclePositionUrl: vehiclePositionURL,
	}
}
