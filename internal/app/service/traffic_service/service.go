package traffic_service

import (
	"log"
	"where-do-i-sit/internal/app"
	serror "where-do-i-sit/internal/app/error"
	"where-do-i-sit/internal/app/sk_api"
	"where-do-i-sit/internal/app/storage"
	"where-do-i-sit/pkg/cache"
)

type TrafficService struct {
	cache cache.Cache
}

func New() TrafficService {
	return TrafficService{
		storage.MemCache,
	}
}

func (t TrafficService) GetStationList() ([]app.Station, error) {
	return sk_api.GetStationList()
}

func (t TrafficService) GetStationByName(s string) (station app.Station, err error) {
	var stations []app.Station
	res, exists := t.cache.Get("stationList")
	if !exists {
		stations, err = sk_api.GetStationList()
		if err != nil {
			return
		}
	} else {
		var ok bool
		stations, ok = res.([]app.Station)
		if !ok {
			log.Println("assertion failed on stationList")
			err = serror.ErrInternal
			return
		}
	}

	for _, st := range stations {
		if st.Name == s {
			station = st
			return
		}
	}

	err = serror.ErrNoSuchStation
	return
}
