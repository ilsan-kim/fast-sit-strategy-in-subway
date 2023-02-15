package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"where-do-i-sit/internal/app"
	"where-do-i-sit/internal/app/error"
	"where-do-i-sit/internal/app/service/traffic_service"
	"where-do-i-sit/internal/app/storage"
	"where-do-i-sit/internal/runtime_util"
	"where-do-i-sit/pkg/cache"
)

type Server struct {
	httpServer     *http.Server
	trafficService traffic_service.TrafficService
	cache          cache.Cache

	IsRunning bool
}

func (s Server) ListenAndServe() error {
	s.startedString()
	return s.httpServer.ListenAndServe()
}

func (s Server) startedString() {
	txt := "\n _____                                 _____  _                _             _ \n/  ___|                               /  ___|| |              | |           | |\n\\ `--.   ___  _ __ __   __ ___  _ __  \\ `--. | |_  __ _  _ __ | |_  ___   __| |\n `--. \\ / _ \\| '__|\\ \\ / // _ \\| '__|  `--. \\| __|/ _` || '__|| __|/ _ \\ / _` |\n/\\__/ /|  __/| |    \\ V /|  __/| |    /\\__/ /| |_| (_| || |   | |_|  __/| (_| |\n\\____/  \\___||_|     \\_/  \\___||_|    \\____/  \\__|\\__,_||_|    \\__|\\___| \\__,_|\n                                                                               \n                                                                               "
	log.Println(txt)
}

func New() *Server {
	server := new(Server)
	server.trafficService = traffic_service.New()

	httpServer := &http.Server{
		Addr:              ":8081",
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:       time.Minute,
	}

	server.httpServer = httpServer
	server.cache = storage.MemCache

	return server
}

func (s Server) RegisterHandler() {
	mux := http.NewServeMux()
	mux.Handle("/hello", middlewareContentType(http.HandlerFunc(helloGetHandler)))
	mux.Handle("/err", middlewareContentType(http.HandlerFunc(errorGetHandler)))
	mux.Handle("/stations", middlewareContentType(http.HandlerFunc(s.stationListGetHandler)))

	s.httpServer.Handler = mux
	return
}

func helloGetHandler(w http.ResponseWriter, r *http.Request) {
	hello := &app.Hello{
		Hello: "World",
	}
	resp, _ := json.Marshal(&hello)
	_, _ = w.Write(resp)
	return
}

func errorGetHandler(w http.ResponseWriter, r *http.Request) {
	e := serror.ErrInvalidRequestTime
	respError(w, e)
}

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

func middlewareContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func respError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch v := err.(type) {
	case serror.Error:
		w.WriteHeader(v.HTTPStatusCode)
		data, err := json.Marshal(v)
		if err != nil {
			fmt.Fprint(w, err.Error())
		}
		fmt.Fprint(w, string(data))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}
	log.Println(err.Error())
}

func (s Server) ShutdownGracefully(ctx context.Context) error {
	s.IsRunning = false
	err := s.httpServer.Shutdown(ctx)
	for {
		cnt := len(runtime_util.GracefulShubdownJob)
		if cnt == 0 {
			break
		}
		log.Println("remained job : ", cnt)
		time.Sleep(time.Second * 5)
	}
	return err
}
