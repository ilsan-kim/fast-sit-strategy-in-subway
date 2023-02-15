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
	config.Conf, err = config.Load(filepath.Join("../../config.json"))
	if err != nil {
		log.Panic(err)
	}
}

func TestGetStationList(t *testing.T) {
	stations, err := GetStationList()
	assert.NoError(t, err)

	for _, station := range stations {
		t.Logf("%s %s: code > %s", station.Line, station.Name, station.Code)
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
	tt := time.Date(2023, 1, 27, 5, 30, 0, 0, time.Local)
	_, _ = GetCongestionForCar("144", tt)
}
