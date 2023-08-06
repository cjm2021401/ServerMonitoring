package influx

import (
	"context"
	"custom_sms/agent/model"
	"custom_sms/config"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func InfluxInit() {
	config.DB = influxdb2.NewClient(config.Env.Database.Url, config.Env.Database.Token)
}

func MemoryStatToInflux(memoryStat model.MemoryStat) error {
	writeAPI := config.DB.WriteAPIBlocking("whatap", "sms")
	hostname, _ := os.Hostname()
	p := influxdb2.NewPointWithMeasurement("Memory").
		AddTag("server_name", hostname).
		AddField("mem_used", memoryStat.MemUsed).
		AddField("swap_used", memoryStat.SwapUsed).
		AddField("mem_total", memoryStat.MemTotal).
		AddField("mem_free", memoryStat.MemFree).
		AddField("mem_available", memoryStat.MemAvailable).
		AddField("buffers", memoryStat.Buffers).
		AddField("cached", memoryStat.Cached).
		AddField("swap_total", memoryStat.SwapTotal).
		AddField("swap_free", memoryStat.SwapFree).
		AddField("mem_shard", memoryStat.Shmem).
		SetTime(time.Now())
	err := writeAPI.WritePoint(context.Background(), p)
	return err
}

func CPUStatToInflux(cpuStats []model.CPUStat) error {
	writeAPI := config.DB.WriteAPIBlocking("whatap", "sms")
	hostname, _ := os.Hostname()
	for _, cpuStat := range cpuStats {
		device := hostname + "_" + cpuStat.Device
		p := influxdb2.NewPointWithMeasurement("CPU").
			AddTag("server_cpu", device).
			AddField("user", cpuStat.User).
			AddField("nice", cpuStat.Nice).
			AddField("system", cpuStat.System).
			AddField("idle", cpuStat.Idle).
			AddField("wait", cpuStat.Wait).
			AddField("irq", cpuStat.Irq).
			AddField("soft_irq", cpuStat.SoftIrq).
			AddField("cpu_used", cpuStat.Used).
			AddField("load_avg1", cpuStat.LoadAvg1).
			AddField("load_avg5", cpuStat.LoadAvg5).
			AddField("load_avg15", cpuStat.LoadAvg15).
			AddField("ctxt", cpuStat.Ctxt).
			AddField("btime", cpuStat.Btime).
			AddField("new_process", cpuStat.NewProcsByForkClone).
			AddField("blocked_process", cpuStat.ProcsBlocked).
			AddField("running_process", cpuStat.ProcsRunning).
			SetTime(time.Now())
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}
	}
	return nil
}
