package model

type MemoryStat struct {
	Timestamp string
	//raw data
	MemTotal     int64
	MemFree      int64
	MemAvailable int64
	Buffers      int64
	Cached       int64
	SwapTotal    int64
	SwapFree     int64
	Shmem        int64

	//percent
	MemUsed  float32
	SwapUsed float32
}
