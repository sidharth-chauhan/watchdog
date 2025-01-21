# watchdog

[![Coverage Status](https://coveralls.io/repos/github/OneBusAway/watchdog/badge.svg?branch=main)](https://coveralls.io/github/OneBusAway/watchdog?branch=main)

Golang-based Watchdog service for OBA REST API servers

# Requirements

Go 1.23 or higher



## Configuration

The watchdog service can be configured using either:
- A **local JSON configuration file** (`--config-file`).
- A **remote JSON configuration URL** (`--config-url`).

### JSON Configuration Format

The JSON configuration file should contain an array of `ObaServer` objects. Example:

```json
[
  {
    "name": "Test Server",
    "id": 1,
    "oba_base_url": "https://test.example.com",
    "oba_api_key": "test-key",
    "gtfs_url": "https://gtfs.example.com",
    "trip_update_url": "https://trip.example.com",
    "vehicle_position_url": "https://vehicle.example.com",
    "gtfs_rt_api_key": "api-key",
    "gtfs_rt_api_value": "api-value"
  }
]
```

# Running

#### **Using a Local Configuration File**

```bash
go run ./cmd/watchdog/ --config-file /path/to/your/config.json
```

#### **Using a Remote Configuration URL with Authentication**

```bash
go run ./cmd/watchdog/ \
  --config-url http://example.com/config.json \
  --config-auth-user admin \
  --config-auth-pass password
```

# Testing

```
go test ./...
```
