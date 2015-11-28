package main

import (
	"collector"
	"misc"
)

func main() {
	misc.Initialize()

	collector_done := make(chan bool)
	go collector.Collector_routine(collector_done)

	<-collector_done
}
