package sk_api

import (
	"github.com/stretchr/testify/assert"
	"log"
	"path/filepath"
	"testing"
	"time"
	"where-do-i-sit/config"
	"where-do-i-sit/internal/app/error"
)

func init() {
	var err error
	log.Println("load configs..")
	config.Conf, err = config.Load(filepath.Join("../../../config.json"))
	if err != nil {
		log.Panic(err)
	}
}

func TestGetStations(t *testing.T) {
	stations, err := GetStations()
	assert.NoError(t, err)
	for k, v := range stations {
		for _, s := range v {
			t.Logf("%s %s: code > %s", s.Line, k, s.Code)
		}
	}
}

func TestGetDow(t *testing.T) {
	sun := time.Date(2023, 1, 22, 23, 0, 0, 0, time.Local)

	for i := 0; i <= 6; i++ {
		tm := sun.AddDate(0, 0, i)
		var expected string
		switch i {
		case 0:
			expected = "SUN"
		case 1:
			expected = "MON"
		case 2:
			expected = "TUE"
		case 3:
			expected = "WED"
		case 4:
			expected = "THU"
		case 5:
			expected = "FRI"
		case 6:
			expected = "SAT"
		}

		assert.Equal(t, expected, getDow(tm))
	}
}

func TestGetHour(t *testing.T) {
	t.Run("5시 이전의 요청은 실패해야 한다.", func(t *testing.T) {
		for i := 0; i <= 23; i++ {
			tt := time.Date(2023, 1, 27, i, 0, 0, 0, time.Local)
			if i <= 5 {
				_, err := getHourIfAvailable(tt)
				assert.Error(t, err, "")
			}
		}
	})

	t.Run("5시 30분 부터 요청 가능, 그 이전은 실패", func(t *testing.T) {
		success := time.Date(2023, 1, 27, 5, 30, 0, 0, time.Local)
		hour, err := getHourIfAvailable(success)
		assert.NoError(t, err, "")
		assert.Equal(t, hour, "05")

		fail := time.Date(2023, 1, 27, 5, 29, 59, 0, time.Local)
		_, err = getHourIfAvailable(fail)
		assert.Error(t, err, "")
		assert.Equal(t, err, serror.ErrInvalidRequestTime)
	})

	t.Run("23시 50분까지 요청 가능, 그 이후는 실패", func(t *testing.T) {
		success := time.Date(2023, 1, 27, 23, 50, 0, 0, time.Local)
		hour, err := getHourIfAvailable(success)
		assert.NoError(t, err, "")
		assert.Equal(t, hour, "23")

		fail := time.Date(2023, 1, 27, 23, 50, 0, 1, time.Local)
		_, err = getHourIfAvailable(fail)
		assert.Error(t, err, "")
		assert.Equal(t, err, serror.ErrInvalidRequestTime)
	})
}

func TestGetCongestionForCar(t *testing.T) {
	// TODO 9시 20분 50초 > 9시 30분의 통계 정보를 보여주는게 맞는지..
	t.Run("정상적인 시간대에 응답이 제데로 내려오는지, 9시 20분 50초의 요청은 9시 30분의 통계 정보를 보여주는지", func(t *testing.T) {
		tt := time.Date(2023, 1, 27, 9, 20, 50, 0, time.Local)
		ret, err := GetStatisticCongestion("216", "217", tt)
		assert.NoError(t, err)
		assert.NotEqual(t, len(ret), 0)
		assert.Equal(t, ret[0].ResponseTime.Hour(), 9)
		assert.Equal(t, ret[0].ResponseTime.Minute(), 30)
	})

	t.Run("비정상적인 시간대에 에러가 잘 오는지", func(t *testing.T) {
		tt := time.Date(2023, 1, 27, 4, 10, 59, 0, time.Local)
		ret, err := GetStatisticCongestion("216", "217", tt)
		assert.Equal(t, err, serror.ErrInvalidRequestTime)
		assert.Equal(t, len(ret), 0)
	})

	t.Run("데이터가 없을때 에러가 잘 오는지", func(t *testing.T) {
		tt := time.Date(2023, 1, 27, 5, 30, 59, 0, time.Local)
		ret, err := GetStatisticCongestion("P555", "P554", tt)
		assert.Equal(t, err, serror.ErrNoData)
		assert.Equal(t, len(ret), 0)
	})
}
