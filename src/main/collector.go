package main

import (
	"time"
	"github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/load"
    "github.com/shirou/gopsutil/mem"
)

const (
	COLLECT_INTERVAL time.Duration = time.Second * 4
)

func collector_routine(done chan bool){
	ticker := updateTicker(COLLECT_INTERVAL, "collector")
    for Running {
        <-ticker.C
        collector_task()
        ticker = updateTicker(COLLECT_INTERVAL, "collector")
    }
    done <- true
}

func collector_task(){
	conn := openConn()
	if conn==nil {
		Err("Local db not found")
	}else{
		
		// por ahora todo va aqui, en el proximo prototipo habra archivos separados
		
		
		now := time.Now()
		unixtime := now.Unix()
		//-=-=-=-=-=-=- Start of collector cycle
		
		// collect CPUTimes
		cpuTimes, _ := cpu.CPUTimes(true)
		for _, cputime := range cpuTimes {
			cpuid := cputime.CPU
			user := cputime.User
			syst := cputime.System
			idle := cputime.Idle
			query := "INSERT INTO "+CPUTIMES_TB+" (cpuid,unixtime,user,sys,idle) values(?,?,?,?,?)"
			err := conn.Exec(query, cpuid, unixtime, user, syst, idle)
			if err != nil {
				Err("saving data to "+CPUTIMES_TB+" table: "+err.Error())
			}
		}
		
		// collect LoadAVG
		loadAvg, _ := load.LoadAvg()
		load1 := loadAvg.Load1
		load5 := loadAvg.Load5
		load15 := loadAvg.Load15
		query := "INSERT INTO "+LOADAVG_TB+" (unixtime,load1,load5,load15) values(?,?,?,?)"
		err := conn.Exec(query, unixtime, load1, load5, load15)
		if err != nil {
			Err("saving data to "+LOADAVG_TB+" table: "+err.Error())
		}
		
		// collect Memory
		ram, _ := mem.VirtualMemory()
		total_ram := int64(ram.Total)
		ram_free := int64(ram.Free)
		ram_used_perc := ram.UsedPercent
		
		// collect swap mem
		swap,_  := mem.SwapMemory()
		swap_total := int64(swap.Total)
		swap_free := int64(swap.Free)
		swap_used_perc := swap.UsedPercent
		
		query = "INSERT INTO "+MEMORY_TB+
		   " (unixtime,total_ram,ram_free,ram_used_percent,total_swap,swap_free,swap_used_percent ) "+
		   "values(?,?,?,?,?,?,?)"
		err = conn.Exec(query, unixtime, total_ram, ram_free, ram_used_perc, swap_total, swap_free, swap_used_perc)
		
		if err != nil {
			Err("saving data to "+MEMORY_TB+" table: "+err.Error())
		}  
		//-=-=-=-=-=-=- End of collector cycle 
		conn.Close()
	}
}