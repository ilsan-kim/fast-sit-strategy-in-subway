package scheduler

import (
	"log"
	"time"
	"where-do-i-sit/internal/app"
	sk_api "where-do-i-sit/internal/app/sk-api"
	"where-do-i-sit/internal/app/storage"
	"where-do-i-sit/pkg/cache"
	"where-do-i-sit/pkg/cache/memcache"
)

type Scheduler struct {
	trafficService app.TrafficService
	cache          cache.Cache

	isRunning bool
}

func (s *Scheduler) InitScheduleJobs() {
	go s.GetStationScheduler(30 * time.Minute)
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		trafficService: sk_api.NewService(),
		cache:          memcache.NewMemCache(),
	}
}

func (s *Scheduler) GetStationScheduler(d time.Duration) {
	station := storage.StationList
	var err error
	for {
		station, err = s.trafficService.GetStationList()
		if err != nil {
			log.Println(err)
		}
		s.cache.Set("stationList", station, d)
		log.Printf("total station %d refreshed\n", len(station))
		time.Sleep(d)
	}
}
