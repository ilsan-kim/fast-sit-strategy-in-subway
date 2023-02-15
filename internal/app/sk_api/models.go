package sk_api

type getStationListResp struct {
	Contents []struct {
		SubwayLine  string `json:"subwayLine"`
		StationName string `json:"stationName"`
		StationCode string `json:"stationCode"`
	} `json:"contents"`
}

type getCongestionForCarResp struct {
	Contents struct {
		SubwayLine  string `json:"subwayLine"`
		StationName string `json:"stationName"`
		StationCode string `json:"stationCode"`
		Stat        []struct {
			PrevStationCode string `json:"prevStationCode"`
			PrevStationName string `json:"prevStationName"`
			UpdnLine        int    `json:"updnLine"`
			DirectAt        int    `json:"directAt"`
			Data            []struct {
				Dow           string `json:"dow"`
				Hh            string `json:"hh"`
				Mm            string `json:"mm"`
				CongestionCar []int  `json:"congestionCar"`
			} `json:"data"`
		} `json:"stat"`
	} `json:"contents"`
}
