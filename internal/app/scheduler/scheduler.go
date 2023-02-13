package scheduler

import (
	"log"
	"time"
	"where-do-i-sit/internal/app"
	sk_api "where-do-i-sit/internal/app/sk-api"
	"where-do-i-sit/internal/app/storage"
	"where-do-i-sit/pkg/cache"
)

type Scheduler struct {
	trafficService app.TrafficService
	cache          cache.Cache

	isRunning bool
}

func (s *Scheduler) InitScheduleJobs() {
	go s.GetStationScheduler(24 * time.Hour)
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		trafficService: sk_api.NewService(),
		cache:          storage.MemCache,
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
