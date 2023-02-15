package traffic_service

import "where-do-i-sit/internal/app"

type TrafficServiceAPI interface {
	GetStationList() ([]app.Station, error)
	GetStationByName(string) (app.Station, error)
}
