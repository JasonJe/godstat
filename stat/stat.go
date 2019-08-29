package stat

import (
	"time"
	"runtime"

	utils "../utils"
	cpu "../cpu"
	memory "../memory"
)

type SysStat struct {
	DateTime utils.FormatTime     `json:"datetime"`
	CpuArray []cpu.CpuStat        `json:"cpuList"`
	memory.MemoryStat
}

func (sysStat *SysStat) CpuUtilization(t int) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(t))
	cpusStat, _ := cpu.CpuTicker()
	<- ticker.C
	cpusStat2, _ := cpu.CpuTicker()

	for i := 0; i < runtime.NumCPU() + 1; i++ {
		cpuStat := cpu.CpuStat{}

		cpuName := cpusStat[i].CPU

		user2 := cpusStat2[i].User
		user := cpusStat[i].User
		nice2 := cpusStat2[i].Nice
		nice := cpusStat[i].Nice
		system2 := cpusStat2[i].System
		system := cpusStat[i].System
		idle2 := cpusStat2[i].Idle
		idle := cpusStat[i].Idle
		iowait2 := cpusStat2[i].Iowait
		iowait := cpusStat[i].Iowait
		irq2 := cpusStat2[i].Irq
		irq := cpusStat[i].Irq
		softirq2 := cpusStat2[i].Softirq
		softirq := cpusStat[i].Softirq
		steal2 := cpusStat2[i].Steal
		steal := cpusStat[i].Steal
		guest2 := cpusStat2[i].Guest
		guest := cpusStat[i].Guest
		guestNice2 := cpusStat2[i].GuestNice
		guestNice := cpusStat[i].GuestNice
		stolen2 := cpusStat2[i].Stolen
		stolen := cpusStat[i].Stolen
		cpu2 := user2 + nice2 + system2 + idle2 + iowait2 + irq2 + softirq2 + steal2 + guest2 + guestNice2 + stolen2
		cpu := user + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice + stolen

		cpuStat.CPU = cpuName
		cpuStat.User = 100 * (user2 - user) / (cpu2 - cpu)
		cpuStat.System = 100 * (system2 - system) / (cpu2 - cpu)
		cpuStat.Idle = 100 * (idle2 - idle) / (cpu2 - cpu)
		cpuStat.Nice = 100 * (nice2 - nice) / (cpu2 - cpu)
		cpuStat.Iowait = 100 * (iowait2 - iowait) / (cpu2 - cpu)
		cpuStat.Irq = 100 * (irq2 - irq) / (cpu2 - cpu)
		cpuStat.Softirq = 100 * (softirq2 - softirq) / (cpu2 - cpu)
		cpuStat.Stolen = 100 * (stolen2 - stolen) / (cpu2 - cpu)
		cpuStat.Steal = 100 * (steal2 - steal) / (cpu2 - cpu)
		cpuStat.Guest = 100 * (guest2 - guest) / (cpu2 - cpu)
		cpuStat.GuestNice = 100 * (guestNice2 - guestNice) / (cpu2 - cpu)

		sysStat.CpuArray = append(sysStat.CpuArray, cpuStat)
	}
}

func (sysStat *SysStat) MemoryInfo() {
	sysStat.MemoryTicker()	
}