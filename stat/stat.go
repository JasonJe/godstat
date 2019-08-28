package stat

import (
	"runtime"

	cpu "../cpu"
	utils "../utils"
)

type SysStat struct {
	CPU       string
	User      float64
	System    float64
	Idle      float64
	Nice      float64
	Iowait    float64
	Irq       float64
	Softirq   float64
	Steal     float64
	Guest     float64
	GuestNice float64
	Stolen    float64
}

func (sysStat *SysStat) CpuTimes() ([]cpu.CpuStat, error) {
	cpusStat := []cpu.CpuStat{}
	filename := "/proc/stat"
	lines, err := utils.ReadLines(filename)
	if err != nil {
		return nil, err
	}

	for i := 0; i < runtime.NumCPU() + 1; i++ {
		cpuStat := cpu.CpuStat{}
		err := cpuStat.GetCpuTimes(lines[i])
		if err != nil {
			return nil, err
		}
		cpusStat = append(cpusStat, cpuStat)
	}
	return cpusStat, nil
}