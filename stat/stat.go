package stat

import (
    "fmt"
	"os"
	"sync"
	"time"
	"runtime"
    "github.com/gosuri/uilive"

	utils  "godstat/utils"
	cpu    "godstat/cpu"
	memory "godstat/memory"
	page   "godstat/page"
	disk   "godstat/disk"
	net    "godstat/net"
	load   "godstat/load"
	swap   "godstat/swap"
	system "godstat/system"
)

type SysStat struct {
	DateTime utils.FormatTime     `json:"datetime"`
	CpuArray []cpu.CpuStat        `json:"cpuList"`
	memory.MemoryStat
	page.PageStat
    DiskList []disk.DiskStat      `json:"diskList"`
    NetList  []net.NetStat        `json:"netList"`
    load.LoadStat
    SwapList []swap.SwapStat      `json:"swapList"`
    system.SystemStat
}

func (sysStat *SysStat) CpuUtilization(t int, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(t))
	cpusStat, _ := cpu.CpuTicker()
	<- ticker.C
	cpusStat2, _ := cpu.CpuTicker()

    (*sysStat).CpuArray = []cpu.CpuStat{}
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

		(*sysStat).CpuArray  = append((*sysStat).CpuArray, cpuStat)
	}
	wg.Done()
}

func (sysStat *SysStat) MemoryInfo(wg *sync.WaitGroup) {
	sysStat.MemoryTicker()
	wg.Done()
}

func (sysStat *SysStat) Paging(t int, wg *sync.WaitGroup) {
	ticker    := time.NewTicker(time.Millisecond * time.Duration(t))
	pageStat  := page.PageStat{}
	pageStat.PageTicker()
	<- ticker.C
	pageStat2 := page.PageStat{}
	pageStat2.PageTicker()

	sysStat.PageIn  = (pageStat2.PageIn  - pageStat.PageIn)  * int64(os.Getpagesize()) * 1
	sysStat.PageOut = (pageStat2.PageOut - pageStat.PageOut) * int64(os.Getpagesize()) * 1

	wg.Done()
}

func (sysStat *SysStat) Disk(t int, wg *sync.WaitGroup) {
    ticker         := time.NewTicker(time.Millisecond * time.Duration(t))
    diskList, err  := disk.DiskTicker()
    if err != nil {
        panic(err)
    }
    <- ticker.C
    diskList2, err := disk.DiskTicker()
    if err != nil {
        panic(err)
    }

    (*sysStat).DiskList = []disk.DiskStat{}
    for index := 0; index < len(diskList); index ++ {
        diskStat := disk.DiskStat{}
        diskStat.Name  = diskList[index].Name
        diskStat.Read  = (diskList2[index].Read  - diskList[index].Read) * 512.0
        diskStat.Write = (diskList2[index].Write - diskList[index].Write) * 512.0

        (*sysStat).DiskList = append((*sysStat).DiskList, diskStat)
    }

    wg.Done()
}

func (sysStat *SysStat) Net(t int, wg *sync.WaitGroup) {
    netDevices, _ := utils.NetDev()

    ticker        := time.NewTicker(time.Millisecond * time.Duration(t))
    netList,  _   := net.NetTicker()
    totalStat1    := netList[len(netList) - 1]
    <- ticker.C
    netList2, _   := net.NetTicker()
    totalStat2    := netList2[len(netList2) - 1]

    (*sysStat).NetList = []net.NetStat{}
    for index, netDev := range netDevices {
        netStat := net.NetStat{}
        netStat.Name = netDev
        netStat.Recv = netList2[index].Recv - netList[index].Recv
        netStat.Send = netList2[index].Send - netList[index].Send

        (*sysStat).NetList = append((*sysStat).NetList, netStat)
    }

    totalStat := net.NetStat{}
    totalStat.Name = "total"
    totalStat.Recv = totalStat2.Recv - totalStat1.Recv
    totalStat.Send = totalStat2.Send - totalStat1.Send

    (*sysStat).NetList = append((*sysStat).NetList, totalStat)
    wg.Done()
}

func (sysStat *SysStat) LoadAvg(wg *sync.WaitGroup) {
    sysStat.LoadTicker()
    wg.Done()
}

func (sysStat *SysStat) Swap(wg *sync.WaitGroup) {
    swapList, _ := swap.SwapTicker()
    (*sysStat).SwapList = swapList
    wg.Done()
}

func (sysStat *SysStat) System(t int, wg *sync.WaitGroup) {
    ticker      := time.NewTicker(time.Millisecond * time.Duration(t))
    systemStat  := system.SystemStat{}
    systemStat.SystemTicker()
    <- ticker.C
    systemStat2 := system.SystemStat{}
    systemStat2.SystemTicker()

    sysStat.Interrupt     = systemStat2.Interrupt     - systemStat.Interrupt
    sysStat.ContextSwitch = systemStat2.ContextSwitch - systemStat.ContextSwitch

    wg.Done()
}

func (sysStat *SysStat) Run(t int) {
    writer        := uilive.New()
    writer.Start()

    fmt.Printf("| --- datetime ---- | ---  cpu(%%) --- | --------- memory usage ---------- | --- paging ---- | - disk total -- | -- net total -- | ---- load avg ---- | ---- swap ----- | --- system ---- |\n")
    fmt.Printf("|           datetime| user|  sys| idel|    used|    free| buffers|  cached|      in|     out|      in|     out|    recv|    send| load| load5| load15|    used|    free|    intr|     csw|\n")

    for {
        startT := time.Now()
        var wg sync.WaitGroup
        wg.Add(9)

        go func(sysStat *SysStat, wg *sync.WaitGroup) {
            sysStat.DateTime = utils.FormatTime(time.Now())
            wg.Done()
        }(sysStat, &wg)
        go sysStat.CpuUtilization(t, &wg)
        go sysStat.MemoryInfo(&wg)
        go sysStat.Paging(t, &wg)
        go sysStat.Disk(t, &wg)
        go sysStat.Net(t, &wg)
        go sysStat.LoadAvg(&wg)
        go sysStat.Swap(&wg)
        go sysStat.System(t, &wg)

        wg.Wait()

        diskListLength := len((*sysStat).DiskList)
        netListLength  := len((*sysStat).NetList)
        swapListLength := len((*sysStat).SwapList)
        tc := time.Since(startT)
        fmt.Fprintf(writer, "|%s|%5.2f|%5.2f|%5.2f|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%5.2f|%6.2f|%7.2f|%8s|%8s|%8s|%8s|\n",
            time.Time((*sysStat).DateTime).Format("2006-01-02 15:04:05"),
            (*sysStat).CpuArray[0].User, (*sysStat).CpuArray[0].System, (*sysStat).CpuArray[0].Idle,
            utils.ByteCountSI(int64((*sysStat).Used)), utils.ByteCountSI(int64((*sysStat).Free)), utils.ByteCountSI(int64((*sysStat).Buffers)), utils.ByteCountSI(int64((*sysStat).Cached)),
            utils.ByteCountSI((*sysStat).PageIn), utils.ByteCountSI((*sysStat).PageOut),
            utils.ByteCountSI(int64((*sysStat).DiskList[diskListLength - 1].Read)), utils.ByteCountSI(int64((*sysStat).DiskList[diskListLength - 1].Write)),
            utils.ByteCountSI(int64((*sysStat).NetList[netListLength - 1].Recv)), utils.ByteCountSI(int64((*sysStat).NetList[netListLength - 1].Send)),
            (*sysStat).Load1, (*sysStat).Load5, (*sysStat).Load15,
            utils.ByteCountSI(int64((*sysStat).SwapList[swapListLength - 1].Used)), utils.ByteCountSI(int64((*sysStat).SwapList[swapListLength - 1].Free)),
            utils.ByteCountSI(int64(sysStat.Interrupt)), utils.ByteCountSI(int64(sysStat.ContextSwitch)))

        fmt.Fprintf(writer, "time const = %v\n", tc)
    }
    writer.Stop()
}

