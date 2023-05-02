package main

import (
	"context"
	"flag"
	"log"
	"time"
	"where-do-i-sit/server/config"
	"where-do-i-sit/server/internal/app/scheduler"
	"where-do-i-sit/server/internal/app/server"
	"where-do-i-sit/server/internal/runtime_util"
)

var configPath *string

func init() {
	configPath = flag.String("conf", "", "path of the config.json file")
	flag.Parse()

	if *configPath == "" {
		log.Println("mandatory flag conf...")
		return
	}
}

func main() {
	runtime_util.GracefulShubdownJob = make(chan struct{}, 1000)

	var err error
	config.Conf, err = config.Load(*configPath)
	if err != nil {
		log.Println(err)
		return
	}
	srv := server.New()
	srv.RegisterHandler()

	sc := scheduler.NewScheduler()
	sc.InitScheduleJobs()

	shutdown := func() {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		srv.ShutdownGracefully(ctx)
	}

	wait := runtime_util.RegisterSignal(shutdown)

	_ = srv.ListenAndServe()

	<-wait
	log.Println("Server has been gracefully shutdown")
}
