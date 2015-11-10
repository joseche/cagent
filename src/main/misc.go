package main

import (
  "fmt"
  "time"
  "os"
)

func updateTicker(interval time.Duration, name string) *time.Ticker {
    nextTick := time.Now().Add(interval)
    fmt.Println(nextTick, name+" - next tick")
    diff := nextTick.Sub(time.Now())
    return time.NewTicker(diff)
}

func File_exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func Info(msg string){
	fmt.Println("info: "+msg)
}

func Err(msg string){
	fmt.Println("Error: "+msg)
}

func Debug(msg string){
	if DEBUG {
		fmt.Println("debug: "+msg)
  	}
}