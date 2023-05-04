package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"where-do-i-sit/client/config"
)

type Station struct {
	Name string
	Line string
	Code string
}

type Stations map[string][]Station

func GetStations() (ret Stations, err error) {
	resp, err := http.Get(config.Conf.API.Url + "stations")
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("failed to decode response : %v", err)
	}
	return
}