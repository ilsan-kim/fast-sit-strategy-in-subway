package traffic_service

import (
	"where-do-i-sit/internal/app"
)

//go:generate mockgen -source=./apis.go  -destination=./mock_apis/mock_apis.go
type TrafficServiceAPI interface {
	GetStationList() ([]app.Station, error)
	GetStationByName(string) (app.Station, error)
}
