package app

import "time"

type Hello struct {
	Hello string `json:"hello"`
}

type Station struct {
	Name string
	Line string
	Code string
}

type Congestion struct {
	From         Station
	ForwardFor   Station
	Congestion   []int
	ResponseTime time.Time
	IsRealtime   bool
}
