package memory

import (
	"bufio"
	"custom_sms/agent/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetMemoryStat(memoryStatChannel chan model.MemoryStat, memoryError chan error) {
	log.Println("Start Collecting Memory metric")
	memoryStat := model.MemoryStat{}
	proc_meminfo, err := os.Open("/proc/meminfo")
	if err != nil {
		memoryError <- err
	}
	defer proc_meminfo.Close()
	timestamp := time.Now()
	memoryStat.Timestamp = timestamp.Format("2006-01-02 15:04:05")
	scan := bufio.NewScanner(proc_meminfo)
	for scan.Scan() {
		data := strings.Fields(scan.Text())
		if strings.HasPrefix(data[0], "MemTotal") {
			memoryStat.MemTotal, _ = strconv.ParseInt(data[1], 10, 64)
		} else if strings.HasPrefix(data[0], "MemFree") {
			memoryStat.MemFree, _ = strconv.ParseInt(data[1], 10, 64)
		} else if strings.HasPrefix(data[0], "MemAvailable") {
			memoryStat.MemAvailable, _ = strconv.ParseInt(data[1], 10, 64)
		} else if strings.HasPrefix(data[0], "Buffers") {
			memoryStat.Buffers, _ = strconv.ParseInt(data[1], 10, 64)
		} else if strings.HasPrefix(data[0], "Cached") {
			memoryStat.Cached, _ = strconv.ParseInt(data[1], 10, 64)
		} else if strings.HasPrefix(data[0], "SwapTotal") {
			memoryStat.SwapTotal, _ = strconv.ParseInt(data[1], 10, 64)
		} else if strings.HasPrefix(data[0], "SwapFree") {
			memoryStat.SwapFree, _ = strconv.ParseInt(data[1], 10, 64)
		} else if strings.HasPrefix(data[0], "Shmem") {
			memoryStat.Shmem, _ = strconv.ParseInt(data[1], 10, 64)
		}
	}
	memoryStat.MemUsed = 100.0 - (float32(memoryStat.MemAvailable) / float32(memoryStat.MemTotal) * 100)
	if memoryStat.SwapTotal != 0 {
		memoryStat.SwapUsed = float32(memoryStat.SwapFree) / float32(memoryStat.SwapTotal) * 100
	}

	err = scan.Err()
	if err != nil {
		memoryError <- err
		return
	}
	memoryStatChannel <- memoryStat
}
