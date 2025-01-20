# watchdog

[![Coverage Status](https://coveralls.io/repos/github/OneBusAway/watchdog/badge.svg?branch=main)](https://coveralls.io/github/OneBusAway/watchdog?branch=main)

Golang-based Watchdog service for OBA REST API servers

# Requirements

Go 1.23 or higher

# Running

```
go run ./cmd/api/ \
  -name <server-name> \
  -id <server-id> \
  -base-url <base-url> \
  -api-key <api-key> \
  -gtfs-url <gtfs-url> \
  -trip-update-url <trip-update-url> \
  -vehicle-position-url <vehicle-position-url>
```

# Testing

```
go test ./...
```
