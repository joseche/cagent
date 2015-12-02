package collector

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"misc"
	"time"
)

const (
	COLLECT_INTERVAL time.Duration = time.Second * 3
)

func Collector_routine(done chan bool) {
	ticker := misc.UpdateTicker(COLLECT_INTERVAL, "collector")
	for misc.Running {
		<-ticker.C
		collector_task()
		ticker = misc.UpdateTicker(COLLECT_INTERVAL, "collector")
	}
	done <- true
}

func collector_task() {
	conn := misc.OpenConn()
	if conn == nil {
		misc.Err("Local db not found")
	} else {

		// por ahora todo va aqui, en el proximo prototipo habra archivos separados

		now := time.Now()
		unixtime := now.Unix()
		//-=-=-=-=-=-=- Start of collector cycle

		// collect LoadAVG
		loadAvg, _ := load.LoadAvg()
		load1 := loadAvg.Load1
		load5 := loadAvg.Load5
		load15 := loadAvg.Load15
		query := "INSERT INTO " + misc.LOADAVG_TB + " (hostid,unixtime,load1,load5,load15) values(?,?,?,?,?)"
		err := conn.Exec(query, misc.Host, unixtime, load1, load5, load15)
		if err != nil {
			misc.Err("saving data to " + misc.LOADAVG_TB + " table: " + err.Error())
		}

		// collect CPUTimes
		cpuTimes, _ := cpu.CPUTimes(true)
		for _, cputime := range cpuTimes {
			cpuid := cputime.CPU

			total := cputime.User + cputime.System + cputime.Idle +
				cputime.Nice + cputime.Iowait + cputime.Irq +
				cputime.Softirq + cputime.Steal + cputime.Guest +
				cputime.GuestNice + cputime.Stolen

			user := cputime.User / total
			syst := cputime.System / total
			idle := cputime.Idle / total
			nice := cputime.Nice / total
			iowait := cputime.Iowait / total
			irq := cputime.Irq / total
			softirq := cputime.Softirq / total
			steal := cputime.Steal / total
			guest := cputime.Guest / total
			guest_nice := cputime.GuestNice / total
			stolen := cputime.Stolen / total

			query := "INSERT INTO " + misc.CPUTIMES_TB +
				" (cpuid,hostid,unixtime,user,sys,idle,nice," +
				"iowait,irq,softirq,steal,guest,guest_nice,stolen) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			err := conn.Exec(query, cpuid, misc.Host, unixtime, user, syst, idle, nice, iowait, irq, softirq, steal, guest, guest_nice, stolen)
			if err != nil {
				misc.Err("saving data to " + misc.CPUTIMES_TB + " table: " + err.Error())
			}
		}

		// collect Memory
		vmem, _ := mem.VirtualMemory()
		vm_total := int64(vmem.Total)
		vm_avail := int64(vmem.Available)
		vm_used := int64(vmem.Used)
		vm_usedpercent := float64(vmem.UsedPercent)
		vm_free := int64(vmem.Free)
		vm_active := int64(vmem.Active)
		vm_inactive := int64(vmem.Inactive)
		vm_buffers := int64(vmem.Buffers)
		vm_cached := int64(vmem.Cached)
		vm_wired := int64(vmem.Wired)
		vm_shared := int64(vmem.Shared)
		query = "INSERT INTO "+misc.VIRTUAL_MEMORY_TB+"(hostid,unixtime,total,available,used,usedpercent,free,active,inactive,buffers,cached,wired,shared) "+
				"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
		err = conn.Exec(query,misc.Host,unixtime,vm_total,vm_avail,vm_used,vm_usedpercent,vm_free,vm_active,vm_inactive,vm_buffers,vm_cached,vm_wired,vm_shared)
		if err != nil {
			misc.Err("saving data to " + misc.VIRTUAL_MEMORY_TB + " table: " + err.Error())
		}
		
		// collect swap mem
		swap, _ := mem.SwapMemory()
		swap_total := int64(swap.Total)
		swap_used := int64(swap.Used)
		swap_free := int64(swap.Free)
		swap_usedperc := swap.UsedPercent
		swap_sin := int64(swap.Sin)
		swap_sout := int64(swap.Sout)
		
		query = "INSERT INTO " + misc.SWAP_MEMORY_TB +
			" (hostid,unixtime,total,used,free,usedpercent,sin,sout ) " +
			"values(?,?,?,?,?,?,?,?)"
		err = conn.Exec(query, misc.Host, unixtime, swap_total, swap_used, swap_free, swap_usedperc, swap_sin, swap_sout)
		if err != nil {
			misc.Err("saving data to " + misc.SWAP_MEMORY_TB + " table: " + err.Error())
		}
		//-=-=-=-=-=-=- End of collector cycle
		conn.Close()
	}
}
