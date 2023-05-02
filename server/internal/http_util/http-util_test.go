package http_util

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetHeader(t *testing.T) {
	headers := map[string]string{
		"key": "value",
	}

	url := "localhost:8080"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	setHeader(headers, req)

	assert.Equal(t, headers["key"], req.Header.Get("key"))
	assert.NotEqual(t, "value2", req.Header.Get("key"))
}

func TestGetAsJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(map[string]string{"hello": "world"})
		if err != nil {
			t.Fatal(err)
		}
		_, err = w.Write(data)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	res, err := GetAsJSON(ts.URL, nil)
	if err != nil {
		t.Logf("err on GetAsJSON %v", err)
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Logf("err on read response body %v", err)
	}

	assert.Equal(t, `{"hello":"world"}`, string(result))
}

func TestPostAsJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			t.Fatal(err)
		}
		_, err = w.Write(reqBody)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	reqBody, err := json.Marshal(map[string]string{"hello": "world"})
	if err != nil {
		t.Fatal(err)
	}

	res, err := PostAsJSON(ts.URL, reqBody, nil)
	if err != nil {
		t.Logf("err on read response body %v", err)
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `{"hello":"world"}`, string(result))
}
