package server

import "watchdog.onebusaway.org/internal/models"

// Config Holds all the configuration settings for our application
type Config struct {
	Port    int
	Env     string
	Servers []models.ObaServer
}

// NewConfig creates a new instance of a Config struct.
func NewConfig(port int, env string, servers []models.ObaServer) *Config {
	return &Config{
		Port:    port,
		Env:     env,
		Servers: servers,
	}
}
