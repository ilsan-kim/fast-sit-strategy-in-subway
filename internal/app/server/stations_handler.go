package server

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s Server) stationListGetHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	res, exists := s.cache.Get("stationList")
	if !exists {
		res, err = s.trafficService.GetStationList()
		if err != nil {
			respError(w, err)
			return
		}
	}
	resp, err := json.Marshal(res)
	if err != nil {
		respError(w, err)
	}
	_, _ = w.Write(resp)
	return
}

func (s Server) stationGetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO 용어사전 만들어야함
	q := r.URL.Query().Get("q")
	res, err := s.trafficService.GetStationByName(q)
	if err != nil {
		respError(w, err)
	}
	resp, err := json.Marshal(res)
	if err != nil {
		respError(w, err)
	}

	_, _ = w.Write(resp)
	return
}

func (s Server) statisticCongestionHandler(w http.ResponseWriter, r *http.Request) {
	stationQ := r.URL.Query().Get("station")
	prevStationQ := r.URL.Query().Get("prevStation")

	// TODO hashmap 으로.. 변경
	station, err := s.trafficService.GetStationByName(stationQ)
	if err != nil {
		respError(w, err)
	}
	prevStation, err := s.trafficService.GetStationByName(prevStationQ)
	if err != nil {
		respError(w, err)
	}

	res, err := s.trafficService.GetStatisticCongestion(station.Code, prevStation.Code, time.Now())
	if err != nil {
		respError(w, err)
	}

	resp, err := json.Marshal(res)
	if err != nil {
		respError(w, err)
	}

	_, _ = w.Write(resp)
	return
}
