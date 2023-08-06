package cpu

import (
	"bufio"
	"custom_sms/agent/model"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetCPUStatAsync(past []model.ProcStatRaw, present chan []model.ProcStatRaw, cpuStats chan []model.CPUStat, cpuError chan error) {
	cpuCoreCount := runtime.NumCPU()
	procStatRawsData := make([]model.ProcStatRaw, cpuCoreCount+1)
	procLoadAvgData := model.ProcLoadAvg{}

	procstatraws := make(chan []model.ProcStatRaw)
	procloadavg := make(chan model.ProcLoadAvg)
	procstatraws_err := make(chan error)
	procloadavg_err := make(chan error)

	go ParseProcStat(procstatraws, procstatraws_err)
	go ParseLoadAvg(procloadavg, procloadavg_err)
	check_procStatRaws := false
	check_procLoadAvg := false
	log.Println("Start Collecting Cpu metric")
	for {
		select {
		case procStatRaws := <-procstatraws:
			procStatRawsData = procStatRaws
			check_procStatRaws = true
			if check_procLoadAvg {
				cpuStats <- Caculate(past, procStatRawsData, procLoadAvgData)
				present <- procStatRawsData
				return
			}
		case procLoadAvg := <-procloadavg:
			procLoadAvgData = procLoadAvg
			check_procLoadAvg = true
			if check_procStatRaws {
				cpuStats <- Caculate(past, procStatRawsData, procLoadAvgData)
				present <- procStatRawsData
				return
			}
		case procStatRaws_err := <-procstatraws_err:
			log.Fatal(procStatRaws_err)
			cpuError <- procStatRaws_err

		case procLoadAvg_err := <-procloadavg_err:
			log.Fatal(procLoadAvg_err)
			cpuError <- procLoadAvg_err
		}
	}

}

func ParseProcStat(procstatraws chan []model.ProcStatRaw, cpu_err chan error) {
	proc_stat, err := os.Open("/proc/stat")
	if err != nil {
		cpu_err <- err
		return
	}
	defer proc_stat.Close()

	cpuCoreCount := runtime.NumCPU()
	procStatRaws := make([]model.ProcStatRaw, cpuCoreCount+1)
	num := 0
	scan := bufio.NewScanner(proc_stat)
	var ctxt, btime, newProcsByForkClone, procsRunning, procsBlocked int64
	for scan.Scan() {
		data := strings.Fields(scan.Text())
		if strings.HasPrefix(data[0], "cpu") {
			procStatRaws[num].Device = data[0]
			procStatRaws[num].User, _ = strconv.ParseInt(data[1], 10, 64)
			procStatRaws[num].System, _ = strconv.ParseInt(data[2], 10, 64)
			procStatRaws[num].Nice, _ = strconv.ParseInt(data[3], 10, 64)
			procStatRaws[num].Idle, _ = strconv.ParseInt(data[4], 10, 64)
			procStatRaws[num].Wait, _ = strconv.ParseInt(data[5], 10, 64)
			procStatRaws[num].Irq, _ = strconv.ParseInt(data[6], 10, 64)
			procStatRaws[num].SoftIrq, _ = strconv.ParseInt(data[7], 10, 64)
			procStatRaws[num].TotalTime = procStatRaws[num].User + procStatRaws[num].System + procStatRaws[num].Nice + procStatRaws[num].Idle + procStatRaws[num].Wait + procStatRaws[num].Irq + procStatRaws[num].SoftIrq
			num++
		} else if data[0] == "ctxt" {
			ctxt, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "btime" {
			btime, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "processes" {
			newProcsByForkClone, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "procs_running" {
			procsRunning, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "procs_blocked" {
			procsRunning, _ = strconv.ParseInt(data[1], 10, 64)
		}
	}

	for index, _ := range procStatRaws {
		procStatRaws[index].Ctxt = ctxt
		procStatRaws[index].Btime = btime
		procStatRaws[index].NewProcsByForkClone = newProcsByForkClone
		procStatRaws[index].ProcsRunning = procsRunning
		procStatRaws[index].ProcsBlocked = procsBlocked
	}

	err = scan.Err()
	if err != nil {
		cpu_err <- err
		return
	}
	procstatraws <- procStatRaws
}

func ParseLoadAvg(procloadavg chan model.ProcLoadAvg, cpu_err chan error) {
	procLoadAvg := model.ProcLoadAvg{}
	proc_loagavg, err := os.Open("/proc/loadavg")
	if err != nil {
		cpu_err <- err
		return
	}
	defer proc_loagavg.Close()

	scan := bufio.NewScanner(proc_loagavg)

	for scan.Scan() {
		data := strings.Fields(scan.Text())
		procLoadAvg.LoadAvg1, _ = strconv.ParseFloat(data[0], 64)
		procLoadAvg.LoadAvg5, _ = strconv.ParseFloat(data[1], 64)
		procLoadAvg.LoadAvg15, _ = strconv.ParseFloat(data[2], 64)
	}

	err = scan.Err()
	if err != nil {
		cpu_err <- err
		return
	}
	procloadavg <- procLoadAvg
}

func Caculate(pastProcStatRaws []model.ProcStatRaw, procStatRawsData []model.ProcStatRaw, procLoadAvgData model.ProcLoadAvg) []model.CPUStat {
	cpuCoreCount := runtime.NumCPU()
	cpuStats := make([]model.CPUStat, cpuCoreCount+1)
	timestamp := time.Now()

	for index, _ := range procStatRawsData {
		cpuStats[index].Timestamp = timestamp.Format("2006-01-02 15:04:05")
		cpuStats[index].Device = procStatRawsData[index].Device
		cpuStats[index].Ctxt = procStatRawsData[index].Ctxt
		cpuStats[index].Btime = procStatRawsData[index].Btime
		cpuStats[index].NewProcsByForkClone = procStatRawsData[index].NewProcsByForkClone
		cpuStats[index].ProcsBlocked = procStatRawsData[index].ProcsBlocked
		cpuStats[index].ProcsRunning = procStatRawsData[index].ProcsRunning

		cpuStats[index].LoadAvg1 = float32(procLoadAvgData.LoadAvg1)
		cpuStats[index].LoadAvg5 = float32(procLoadAvgData.LoadAvg5)
		cpuStats[index].LoadAvg15 = float32(procLoadAvgData.LoadAvg15)

		TimeInterver := float32(procStatRawsData[index].TotalTime) - float32(pastProcStatRaws[index].TotalTime)

		cpuStats[index].User = (float32(procStatRawsData[index].User) - float32(pastProcStatRaws[index].User)) / TimeInterver * 100.0
		cpuStats[index].Nice = (float32(procStatRawsData[index].Nice) - float32(pastProcStatRaws[index].Nice)) / TimeInterver * 100.0
		cpuStats[index].System = (float32(procStatRawsData[index].System) - float32(pastProcStatRaws[index].System)) / TimeInterver * 100.0
		cpuStats[index].Idle = (float32(procStatRawsData[index].Idle) - float32(pastProcStatRaws[index].Idle)) / TimeInterver * 100.0
		cpuStats[index].Wait = (float32(procStatRawsData[index].Wait) - float32(pastProcStatRaws[index].Wait)) / TimeInterver * 100.0
		cpuStats[index].Irq = (float32(procStatRawsData[index].Irq) - float32(pastProcStatRaws[index].Irq)) / TimeInterver * 100.0
		cpuStats[index].SoftIrq = (float32(procStatRawsData[index].SoftIrq) - float32(pastProcStatRaws[index].SoftIrq)) / TimeInterver * 100.0
		cpuStats[index].Used = 100.0 - cpuStats[index].Idle
	}

	return cpuStats
}

func ParseProcStatFirstTime() ([]model.ProcStatRaw, error) {
	proc_stat, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer proc_stat.Close()

	cpuCoreCount := runtime.NumCPU()
	procStatRaws := make([]model.ProcStatRaw, cpuCoreCount+1)
	num := 0
	scan := bufio.NewScanner(proc_stat)
	var ctxt, btime, newProcsByForkClone, procsRunning, procsBlocked int64
	for scan.Scan() {
		data := strings.Fields(scan.Text())
		if strings.HasPrefix(data[0], "cpu") {
			procStatRaws[num].Device = data[0]
			procStatRaws[num].User, _ = strconv.ParseInt(data[1], 10, 64)
			procStatRaws[num].System, _ = strconv.ParseInt(data[2], 10, 64)
			procStatRaws[num].Nice, _ = strconv.ParseInt(data[3], 10, 64)
			procStatRaws[num].Idle, _ = strconv.ParseInt(data[4], 10, 64)
			procStatRaws[num].Wait, _ = strconv.ParseInt(data[5], 10, 64)
			procStatRaws[num].Irq, _ = strconv.ParseInt(data[6], 10, 64)
			procStatRaws[num].SoftIrq, _ = strconv.ParseInt(data[7], 10, 64)
			procStatRaws[num].TotalTime = procStatRaws[num].User + procStatRaws[num].System + procStatRaws[num].Nice + procStatRaws[num].Idle + procStatRaws[num].Wait + procStatRaws[num].Irq + procStatRaws[num].SoftIrq
			num++
		} else if data[0] == "ctxt" {
			ctxt, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "btime" {
			btime, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "processes" {
			newProcsByForkClone, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "procs_running" {
			procsRunning, _ = strconv.ParseInt(data[1], 10, 64)
		} else if data[0] == "procs_blocked" {
			procsRunning, _ = strconv.ParseInt(data[1], 10, 64)
		}
	}

	for index, _ := range procStatRaws {
		procStatRaws[index].Ctxt = ctxt
		procStatRaws[index].Btime = btime
		procStatRaws[index].NewProcsByForkClone = newProcsByForkClone
		procStatRaws[index].ProcsRunning = procsRunning
		procStatRaws[index].ProcsBlocked = procsBlocked
	}

	err = scan.Err()
	if err != nil {
		return nil, err
	}
	return procStatRaws, err
}
