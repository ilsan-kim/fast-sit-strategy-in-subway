package sk_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"where-do-i-sit/config"
	"where-do-i-sit/internal/app"
	"where-do-i-sit/internal/app/error"
	http_util "where-do-i-sit/internal/http_util"
)

var dowMap = map[int]string{
	0: "SUN",
	1: "MON",
	2: "TUE",
	3: "WED",
	4: "THU",
	5: "FRI",
	6: "SAT",
}

func getDow(t time.Time) string {
	return dowMap[int(t.Weekday())]
}

func getHourIfAvailable(t time.Time) (string, error) {
	beforeCriteria := time.Date(t.Year(), t.Month(), t.Day(), 5, 30, 0, 0, t.Location())
	afterCriteria := time.Date(t.Year(), t.Month(), t.Day(), 23, 50, 0, 0, t.Location())
	if t.Before(beforeCriteria) || t.After(afterCriteria) {
		return "", serror.ErrInvalidRequestTime
	}
	hour := strconv.Itoa(t.Hour())
	if len(hour) == 1 {
		hour = "0" + hour
	}
	return hour, nil
}

func getDefaultHeader() map[string]string {
	return map[string]string{
		"appkey": config.Conf.SkApi.AppKey,
	}
}

func GetStationList() (ret []app.Station, err error) {
	var res getStationListResp
	url := "https://apis.openapi.sk.com/puzzle/subway/stations"
	headers := getDefaultHeader()

	resp, err := http_util.GetAsJSON(url, headers)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, serror.ErrExternalService
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(respBody, &res)

	for _, station := range res.Contents {
		ret = append(ret, app.Station{
			Name: station.StationName,
			Line: station.SubwayLine,
			Code: station.StationCode,
		})
	}

	return
}

func GetCongestionForCar(stationCode string, time time.Time) (ret any, err error) {
	var res getCongestionForCarResp
	dow := getDow(time)
	hour, err := getHourIfAvailable(time)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://apis.openapi.sk.com/puzzle/congestion-car/stat/stations/%s?dow=%s&hh=%s", stationCode, dow, hour)
	headers := getDefaultHeader()

	resp, err := http_util.GetAsJSON(url, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, serror.ErrExternalService
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(respBody, &res)

	log.Println(res)
	return nil, nil
}
