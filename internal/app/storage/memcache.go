package storage

import (
	"where-do-i-sit/internal/app"
	"where-do-i-sit/pkg/cache"
	"where-do-i-sit/pkg/cache/memcache"
)

func init() {
	MemCache = memcache.NewMemCache()
}

var MemCache cache.Cache
var StationList []app.Station
var Stations app.Stations
