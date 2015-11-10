package main

const (
  DEBUG bool = true
  MASTER_DIR  string = "/opt/cloudmetrics/"

)

var Running bool = true

func main(){
	
	//collector_done := make(chan bool)
	//go collector_routine(collector_done)
	
	generate_ssl()
	
	//<- collector_done
}