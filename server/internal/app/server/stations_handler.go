package server

import (
	"encoding/json"
	"net/http"
	"time"
	"where-do-i-sit/server/internal/app/error"
)

func (s Server) stationListGetHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	res, exists := s.cache.Get("stations")
	if !exists {
		res, err = s.trafficService.GetStations()
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
	lineQ := r.URL.Query().Get("line")
	if lineQ == "" {
		e := serror.ErrNoRequiredParam
		respError(w, e.FormatMsg("line"))
		return
	}

	stationQ := r.URL.Query().Get("station")
	if stationQ == "" {
		e := serror.ErrNoRequiredParam
		respError(w, e.FormatMsg("station"))
		return
	}

	res, err := s.trafficService.GetStationByName(stationQ, lineQ)
	if err != nil {
		respError(w, err)
		return
	}
	resp, err := json.Marshal(res)
	if err != nil {
		respError(w, err)
		return
	}

	_, _ = w.Write(resp)
	return
}

func (s Server) statisticCongestionHandler(w http.ResponseWriter, r *http.Request) {
	lineQ := r.URL.Query().Get("line")
	if lineQ == "" {
		e := serror.ErrNoRequiredParam
		respError(w, e.FormatMsg("line"))
		return
	}

	stationQ := r.URL.Query().Get("station")
	if stationQ == "" {
		e := serror.ErrNoRequiredParam
		respError(w, e.FormatMsg("station"))
		return
	}

	prevStationQ := r.URL.Query().Get("prevStation")
	if prevStationQ == "" {
		e := serror.ErrNoRequiredParam
		respError(w, e.FormatMsg("prevStation"))

		return
	}

	// TODO hashmap 으로.. 변경
	station, err := s.trafficService.GetStationByName(stationQ, lineQ)
	if err != nil {
		respError(w, err)
		return
	}
	prevStation, err := s.trafficService.GetStationByName(prevStationQ, lineQ)
	if err != nil {
		respError(w, err)
		return
	}

	res, err := s.trafficService.GetStatisticCongestion(station.Code, prevStation.Code, time.Now())
	if err != nil {
		respError(w, err)
		return
	}

	resp, err := json.Marshal(res)
	if err != nil {
		respError(w, err)
		return
	}

	_, _ = w.Write(resp)
	return
}
