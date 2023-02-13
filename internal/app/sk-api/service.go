package sk_api

import (
	"where-do-i-sit/internal/app"
)

type SKAPIService struct{}

func (s *SKAPIService) GetStationList() ([]app.Station, error) {
	return GetStationList()
}

func NewService() *SKAPIService {
	return &SKAPIService{}
}
