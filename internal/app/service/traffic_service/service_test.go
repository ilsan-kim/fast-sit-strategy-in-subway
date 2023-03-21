package traffic_service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"path/filepath"
	"testing"
	"time"
	"where-do-i-sit/config"
	"where-do-i-sit/internal/app"
	serror "where-do-i-sit/internal/app/error"
	"where-do-i-sit/internal/app/storage"
	"where-do-i-sit/pkg/cache/mock_cache"
)

func init() {
	var err error
	log.Println("load configs..")
	config.Conf, err = config.Load(filepath.Join("../../../../config.json"))
	if err != nil {
		log.Panic(err)
	}
}

func TestTrafficService_GetStations(t *testing.T) {
	ts := New(storage.MemCache)

	stationList, err := ts.GetStations()
	log.Println(stationList)
	assert.NoError(t, err)
	assert.True(t, len(stationList) > 0)
}

func TestTrafficService_GetStationByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stations := make(app.Stations)
	station := []app.Station{{
		Name: "옥수수역",
		Line: "2호선",
		Code: "CORN",
	}}
	stations["옥수수역"] = station

	t.Run("입력한 역이 캐시에 있는 경우", func(t *testing.T) {
		mc := mock_cache.NewMockCache(ctrl)
		mc.EXPECT().Get("stations").Return(stations, true)
		ts := New(mc)
		st, err := ts.GetStationByName("옥수수역", "2호선")

		assert.NoError(t, err)
		assert.Equal(t, stations["옥수수역"][0].Name, st.Name)
	})

	t.Run("입력한 역이 캐시에 없는 역인 경우", func(t *testing.T) {
		mc := mock_cache.NewMockCache(ctrl)
		mc.EXPECT().Get("stations").Return(stations, true)
		ts := New(mc)
		st, err := ts.GetStationByName("치킨역", "2호선")

		assert.Equal(t, serror.ErrNoSuchStation, err)
		assert.Equal(t, "", st.Name)
	})

	t.Run("캐시가 없는 경우 >> trafficService.GetStationList 에서 잘 받아오는지 확인 // 단일 노선", func(t *testing.T) {
		mc := mock_cache.NewMockCache(ctrl)
		mc.EXPECT().Get("stations").Return([]app.Station{}, false)
		ts := New(mc)
		sl, _ := ts.GetStations()
		mc.EXPECT().Set("stations", sl, time.Hour*24).Return()
		st, err := ts.GetStationByName("신대방역", "2호선")
		assert.NoError(t, err)
		log.Println(st)
		assert.Equal(t, "신대방역", st.Name)
		assert.Equal(t, "2호선", st.Line)
		assert.Equal(t, "231", st.Code)
	})

	t.Run("캐시가 없는 경우 >> trafficService.GetStationList 에서 잘 받아오는지 확인 // 환승 노선", func(t *testing.T) {
		mc := mock_cache.NewMockCache(ctrl)
		mc.EXPECT().Get("stations").Return([]app.Station{}, false).AnyTimes()
		ts := New(mc)
		sl, _ := ts.GetStations()
		mc.EXPECT().Set("stations", sl, time.Hour*24).Return().AnyTimes()
		st, err := ts.GetStationByName("시청역", "2호선")
		assert.NoError(t, err)
		assert.Equal(t, "시청역", st.Name)
		assert.Equal(t, "2호선", st.Line)
		assert.Equal(t, "201", st.Code)

		st2, err := ts.GetStationByName("시청역", "1호선")
		assert.NoError(t, err)
		assert.Equal(t, "시청역", st2.Name)
		assert.Equal(t, "1호선", st2.Line)
		assert.Equal(t, "132", st2.Code)
	})
}
