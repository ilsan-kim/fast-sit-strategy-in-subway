package traffic_service

import (
	"log"
	"time"
	"where-do-i-sit/internal/app"
	serror "where-do-i-sit/internal/app/error"
	"where-do-i-sit/internal/app/sk_api"
	"where-do-i-sit/pkg/cache"
)

type TrafficService struct {
	cache cache.Cache
}

func New(cache cache.Cache) TrafficService {
	return TrafficService{
		cache,
	}
}

func (t TrafficService) GetStations() (app.Stations, error) {
	return sk_api.GetStations()
}

func (t TrafficService) GetStationByName(s string, line string) (station app.Station, err error) {
	var stations app.Stations
	res, exists := t.cache.Get("stations")
	if !exists {
		stations, err = t.GetStations()
		t.cache.Set("stations", stations, 24*time.Hour)
		if err != nil {
			return
		}
	} else {
		var ok bool
		stations, ok = res.(app.Stations)
		if !ok {
			log.Println("assertion failed on stationList")
			err = serror.ErrInternal
			return
		}
	}

	if sts, ok := stations[s]; ok {
		if len(sts) == 0 {
			return sts[0], nil
		} else {
			for _, s := range sts {
				if s.Line == line {
					return s, nil
				}
			}
		}
	} else {
		err = serror.ErrNoSuchStation
	}
	return
}

func (t TrafficService) GetStatisticCongestion(stationCode, prevStationCode string, time time.Time) (ret []app.Congestion, err error) {
	ret, err = sk_api.GetStatisticCongestion(stationCode, prevStationCode, time)
	return
}

func (t TrafficService) GetRealtimeCongestion(stationCode, prevStationCode string) (app.Congestion, error) {
	//TODO implement me
	panic("implement me")
}
