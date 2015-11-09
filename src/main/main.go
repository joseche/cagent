package main

const (
  DEBUG bool = true
)

var Running bool = true

func main(){
	
	collector_done := make(chan bool)
	go collector_routine(collector_done)
	
	
	
	<- collector_done
}