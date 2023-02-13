package app

type TrafficService interface {
	GetStationList() ([]Station, error)
}
