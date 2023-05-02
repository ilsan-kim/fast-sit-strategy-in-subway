package traffic_service

import (
	"time"
	"where-do-i-sit/server/internal/app"
)

//go:generate mockgen -source=./apis.go  -destination=./mock_apis/mock_apis.go
type TrafficServiceAPI interface {
	GetStations() (app.Stations, error)
	GetStationByName(string, string) (app.Station, error)
	GetStatisticCongestion(stationCode, prevStationCode string, time time.Time) ([]app.Congestion, error)
	GetRealtimeCongestion(stationCode, prevStationCode string) (app.Congestion, error)
}
