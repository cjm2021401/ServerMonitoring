package main

import (
	"custom_sms/agent/cpu"
	"custom_sms/agent/memory"
	"custom_sms/agent/model"
	"custom_sms/alert"
	"custom_sms/config"
	"custom_sms/influx"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	err := config.GetEnvironmentVariable()
	if err != nil {
		log.Panic(err)
	}
	influx.InfluxInit()
	wg.Add(2)
	go memoryMonitor()
	go cpuMonitor()
	wg.Wait()

}
func cpuMonitor() {
	procStatRaws, err := cpu.ParseProcStatFirstTime()

	cpuError := make(chan error)
	cpuStats := make(chan []model.CPUStat)
	present := make(chan []model.ProcStatRaw)

	if err != nil {
		log.Println("Cant Collect CPU Data", err)
	}
	for {
		time.Sleep(10 * time.Second)
		log.Println(procStatRaws[0].TotalTime)
		go cpu.GetCPUStatAsync(procStatRaws, present, cpuStats, cpuError)
		for i := 0; i < 2; i++ {
			select {
			case CPUStats := <-cpuStats:
				log.Println("CPU Used Percent(", CPUStats[0].Timestamp, ")", CPUStats[0].Used, "%")
				err := influx.CPUStatToInflux(CPUStats)
				if err != nil {
					log.Fatal("Cant Insert Memory Data", err)
				}
				err = alert.CPUAlert(CPUStats)
				if err != nil {
					log.Println("Cant Send Message to Slack", err)
				}
			case PRESENT := <-present:
				procStatRaws = PRESENT
			case cpuError := <-cpuError:
				log.Println("Cant Collect CPU Data", cpuError)
				break
			}
		}

	}
}

func memoryMonitor() {

	memoryStat := make(chan model.MemoryStat)
	memoryError := make(chan error)
	for {
		go memory.GetMemoryStat(memoryStat, memoryError)

		select {
		case MEMORYStats := <-memoryStat:
			log.Println("Memory Used Percent(", MEMORYStats.Timestamp, ")", MEMORYStats.MemUsed, "%")
			err := influx.MemoryStatToInflux(MEMORYStats)
			if err != nil {
				log.Fatal("Cant Insert Memory Data", err)
			}
			err = alert.MemoryAlert(MEMORYStats)
			if err != nil {
				log.Println("Cant Send Message to Slack", err)
			}

		case memoryError := <-memoryError:
			log.Println("Cant Collect Memory Data", memoryError)
		}
		time.Sleep(10 * time.Second)

	}
}
