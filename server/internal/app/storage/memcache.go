package storage

import (
	"where-do-i-sit/pkg/cache"
	"where-do-i-sit/pkg/cache/memcache"
	"where-do-i-sit/server/internal/app"
)

func init() {
	MemCache = memcache.NewMemCache()
}

var MemCache cache.Cache
var StationList []app.Station
var Stations app.Stations
