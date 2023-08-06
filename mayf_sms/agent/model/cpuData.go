package model

type ProcStatRaw struct {
	Device    string
	TotalTime int64
	User      int64
	System    int64
	Nice      int64
	Idle      int64
	Wait      int64
	Irq       int64
	SoftIrq   int64

	Ctxt  int64
	Btime int64

	NewProcsByForkClone int64
	ProcsRunning        int64
	ProcsBlocked        int64
}

type ProcLoadAvg struct {
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64
}

type CPUStat struct {
	Timestamp string
	Device    string
	//persent
	User    float32
	Nice    float32
	System  float32
	Idle    float32
	Wait    float32
	Irq     float32
	SoftIrq float32
	Used    float32

	LoadAvg1  float32
	LoadAvg5  float32
	LoadAvg15 float32

	//value
	Ctxt                int64
	Btime               int64
	NewProcsByForkClone int64
	ProcsBlocked        int64
	ProcsRunning        int64
}
