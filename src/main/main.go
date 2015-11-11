package main

const (
	DEBUG      bool   = true
	MASTER_DIR string = "/opt/cloudmetrics/"
)

var Running bool = true
var host string

func main() {
	initialize()

	collector_done := make(chan bool)
	go collector_routine(collector_done)

	<-collector_done
}
