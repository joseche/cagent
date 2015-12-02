package main

import (
	"collector"
	"misc"
	"pusher"
)

func main() {
	misc.Initialize()

	collector_done := make(chan bool)
	go collector.Collector_routine(collector_done)
	
	pusher_done := make(chan bool)
	go pusher.Pusher_routine(pusher_done)

	<-collector_done
	<-pusher_done
}
