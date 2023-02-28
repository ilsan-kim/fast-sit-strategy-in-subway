package traffic_service

import (
	"time"
	"where-do-i-sit/internal/app"
)

//go:generate mockgen -source=./apis.go  -destination=./mock_apis/mock_apis.go
type TrafficServiceAPI interface {
	GetStationList() ([]app.Station, error)
	GetStationByName(string) (app.Station, error)
	GetStatisticCongestion(stationCode, prevStationCode string, time time.Time) ([]app.Congestion, error)
	GetRealtimeCongestion(stationCode, prevStationCode string) (app.Congestion, error)
}
