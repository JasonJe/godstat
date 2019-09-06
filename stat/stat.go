package stat

import (
    "fmt"
	"os"
	"time"
	"runtime"

	utils  "../utils"
	cpu    "../cpu"
	memory "../memory"
	page   "../page"
	disk   "../disk"
)

type SysStat struct {
	DateTime utils.FormatTime     `json:"datetime"`
	CpuArray []cpu.CpuStat        `json:"cpuList"`
	memory.MemoryStat
	page.PageStat
    DiskList []disk.DiskStat      `json:"diskList"`
}

func (sysStat *SysStat) CpuUtilization(t int) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(t))
	cpusStat, _ := cpu.CpuTicker()
	<- ticker.C
	cpusStat2, _ := cpu.CpuTicker()

	for i := 0; i < runtime.NumCPU() + 1; i++ {
		cpuStat := cpu.CpuStat{}

		cpuName := cpusStat[i].CPU

		user2      := cpusStat2[i].User
		user       := cpusStat[i].User
		nice2      := cpusStat2[i].Nice
		nice       := cpusStat[i].Nice
		system2    := cpusStat2[i].System
		system     := cpusStat[i].System
		idle2      := cpusStat2[i].Idle
		idle       := cpusStat[i].Idle
		iowait2    := cpusStat2[i].Iowait
		iowait     := cpusStat[i].Iowait
		irq2       := cpusStat2[i].Irq
		irq        := cpusStat[i].Irq
		softirq2   := cpusStat2[i].Softirq
		softirq    := cpusStat[i].Softirq
		steal2     := cpusStat2[i].Steal
		steal      := cpusStat[i].Steal
		guest2     := cpusStat2[i].Guest
		guest      := cpusStat[i].Guest
		guestNice2 := cpusStat2[i].GuestNice
		guestNice  := cpusStat[i].GuestNice
		stolen2    := cpusStat2[i].Stolen
		stolen     := cpusStat[i].Stolen
		cpu2       := user2 + nice2 + system2 + idle2 + iowait2 + irq2 + softirq2 + steal2 + guest2 + guestNice2 + stolen2
		cpu        := user + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice + stolen

		cpuStat.CPU       = cpuName
		cpuStat.User      = 100 * (user2 - user) / (cpu2 - cpu)
		cpuStat.System    = 100 * (system2 - system) / (cpu2 - cpu)
		cpuStat.Idle      = 100 * (idle2 - idle) / (cpu2 - cpu)
		cpuStat.Nice      = 100 * (nice2 - nice) / (cpu2 - cpu)
		cpuStat.Iowait    = 100 * (iowait2 - iowait) / (cpu2 - cpu)
		cpuStat.Irq       = 100 * (irq2 - irq) / (cpu2 - cpu)
		cpuStat.Softirq   = 100 * (softirq2 - softirq) / (cpu2 - cpu)
		cpuStat.Stolen    = 100 * (stolen2 - stolen) / (cpu2 - cpu)
		cpuStat.Steal     = 100 * (steal2 - steal) / (cpu2 - cpu)
		cpuStat.Guest     = 100 * (guest2 - guest) / (cpu2 - cpu)
		cpuStat.GuestNice = 100 * (guestNice2 - guestNice) / (cpu2 - cpu)

		sysStat.CpuArray  = append(sysStat.CpuArray, cpuStat)
	}
}

func (sysStat *SysStat) MemoryInfo() {
	sysStat.MemoryTicker()	
}

func (sysStat *SysStat) Paging(t int) {
	ticker    := time.NewTicker(time.Millisecond * time.Duration(t))
	pageStat  := page.PageStat{}
	pageStat.PageTicker()
	<- ticker.C
	pageStat2 := page.PageStat{}
	pageStat2.PageTicker()

	sysStat.PageIn  = (pageStat2.PageIn  - pageStat.PageIn)  * int64(os.Getpagesize()) * 1
	sysStat.PageOut = (pageStat2.PageOut - pageStat.PageOut) * int64(os.Getpagesize()) * 1
}

func (sysStat *SysStat) Disk(t int, totalDiskStat *disk.DiskStat) {
    blockDevices, _ := utils.GetDiskDev() 

    ticker          := time.NewTicker(time.Millisecond * time.Duration(t))
    diskList,  _    := disk.DiskTicker(totalDiskStat)
    <- ticker.C
    diskList2, _    := disk.DiskTicker(totalDiskStat)

    (*sysStat).DiskList = []disk.DiskStat{} 
    for _, name := range blockDevices {
        diskStat := disk.DiskStat{}
        diskStat.Name    = name
        diskStat.Read    = (diskList2[name].Read  - diskList[name].Read)  * 512.0
        diskStat.Write   = (diskList2[name].Write - diskList[name].Write) * 512.0
        
       (*sysStat).DiskList = append((*sysStat).DiskList, diskStat)
    }
}

func (sysStat *SysStat) Run(t int) {
    totalDiskStat := &disk.DiskStat{"total", 0.0, 0.0}
    for {
        sysStat.CpuUtilization(t)
        sysStat.MemoryInfo()
        sysStat.Paging(t)
        sysStat.Disk(t, totalDiskStat)
        
        time.Sleep(time.Second)
        
        fmt.Println(sysStat.DiskList)
    }
}

