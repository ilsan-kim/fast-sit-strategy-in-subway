package scheduler

import (
	"log"
	"time"
	"where-do-i-sit/internal/app/service/traffic_service"
	"where-do-i-sit/internal/app/storage"
	"where-do-i-sit/pkg/cache"
)

var (
	_ traffic_service.TrafficServiceAPI = (*traffic_service.TrafficService)(nil)
)

type Scheduler struct {
	trafficService traffic_service.TrafficServiceAPI
	cache          cache.Cache
}

func (s *Scheduler) InitScheduleJobs() {
	go s.GetStationScheduler(24 * time.Hour)
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		trafficService: traffic_service.New(storage.MemCache),
		cache:          storage.MemCache,
	}
}

func (s *Scheduler) GetStationScheduler(d time.Duration) {
	stations := storage.Stations
	var err error
	for {
		stations, err = s.trafficService.GetStations()
		if err != nil {
			log.Println(err)
		}
		// TODO 스테이션을 map 형태로 저장하여 검색속도 빠르게 만들기
		s.cache.Set("stations", stations, d)
		log.Printf("total station %d refreshed\n", len(stations))
		time.Sleep(d)
	}
}
