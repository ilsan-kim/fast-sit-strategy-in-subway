package sk_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"where-do-i-sit/server/config"
	"where-do-i-sit/server/internal/app"
	"where-do-i-sit/server/internal/app/error"
	"where-do-i-sit/server/internal/http_util"
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

func GetStations() (ret app.Stations, err error) {
	ret = make(app.Stations)
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
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, serror.ErrQuotaExceed
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(respBody, &res)

	for _, station := range res.Contents {
		if val, ok := ret[station.StationName]; !ok {
			s := app.Station{
				Name: station.StationName,
				Line: station.SubwayLine,
				Code: station.StationCode,
			}
			ret[station.StationName] = append(ret[station.StationName], s)
		} else {
			s := app.Station{
				Name: station.StationName,
				Line: station.SubwayLine,
				Code: station.StationCode,
			}

			ret[station.StationName] = append(val, s)
		}
	}
	return
}

func GetStatisticCongestion(stationCode string, prevStationCode string, t time.Time) (ret []app.Congestion, err error) {
	// 50분이 지났다면, 다음 시간(hour)로 넘겨 계산한다.
	if t.Minute() >= 50 {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, t.Second(), t.Nanosecond(), t.Location())
	}
	var res getStatisticCongestionResp
	dow := getDow(t)
	hour, err := getHourIfAvailable(t)
	if err != nil {
		return ret, err
	}
	url := fmt.Sprintf("https://apis.openapi.sk.com/puzzle/congestion-car/stat/stations/%s?dow=%s&hh=%s", stationCode, dow, hour)
	headers := getDefaultHeader()

	resp, err := http_util.GetAsJSON(url, headers)
	if err != nil {
		log.Println(err)
		return ret, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return ret, serror.ErrExternalService
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return ret, serror.ErrQuotaExceed
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(respBody, &res)

	for _, stat := range res.Contents.Stat {
		if stat.PrevStationCode == prevStationCode {
			//firstMetTime := time.Time{}
			for _, data := range stat.Data {

				// 혼잡도가 응답되지 않은 stat.Data 는 제외한다.
				congestion := 0
				for _, c := range data.CongestionCar {
					congestion += c
				}
				if congestion == 0 {
					continue
				}

				hInt, _ := strconv.Atoi(data.Hh)
				mInt, _ := strconv.Atoi(data.Mm)

				dataTime := makeTime(t, hInt, mInt)
				// 10분 이내의 dataTime 을 가져온다.
				if t.After(dataTime) || t.Add(time.Minute*10).Before(dataTime) {
					continue
				}

				if stat.EndStationCode == "211-R" {
					stat.EndStationName = "순환"
				}

				ret = append(ret, app.Congestion{
					From: app.Station{
						Name: stat.PrevStationName,
						Line: res.Contents.SubwayLine,
						Code: stat.PrevStationCode,
					},
					ForwardFor: app.Station{
						Name: stat.EndStationName,
						Line: res.Contents.SubwayLine,
						Code: stat.EndStationCode,
					},
					Congestion:   data.CongestionCar,
					ResponseTime: dataTime,
					IsRealtime:   false,
				})
				log.Printf("[%s:%s]%s에서 오는 %s 출발 %s행 기차 칸별 혼잡도 %v", data.Hh, data.Mm, stat.PrevStationName, stat.StartStationName, stat.EndStationName, data.CongestionCar)
			}
		}
	}

	if len(ret) == 0 {
		return nil, serror.ErrNoData
	}

	return ret, nil
}

func makeTime(t time.Time, hh, mm int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hh, mm, 0, 0, time.Local)
}

func GetRealtimeCongestion(stationCode string, prevStationCode string, t time.Time) (ret []app.Congestion, err error) {

	return
}
