package stat

import (
    "fmt"
    "os"
    "sync"
    "time"
    "runtime"
    "github.com/gosuri/uilive"

    utils      "godstat/utils"
    core       "godstat/core"
)

type SysStat struct {
    DateTime  utils.FormatTime     `json:"datetime"`
    CpuArray  []core.CpuStat        `json:"cpuList"`
    core.MemoryStat 
    core.PageStat 
    DiskList  []core.DiskStat      `json:"diskList"`
    NetList   []core.NetStat        `json:"netList"`
    core.LoadStat
    SwapList  []core.SwapStat      `json:"swapList"`
    core.SystemStat
    Socket     core.SocketStat
    RawSocket  core.RawSocketStat
    UnixSocket core.UnixSocketStat 
    TCP        core.TCPStat 
    UDP        core.UDPStat
    core.FileSystemStat
    IOList  []core.IOStat            `json:"diskList"`
    AIO        core.AIOStat
    Proc       core.ProcStat
    IPC        core.IPCStat
    Zone       core.ZoneStat 
    Lock       core.LockStat 
    VM         core.VMStat 
}

func (sysStat *SysStat) CpuUtilization(t int, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(t))
	cpusStat, _ := core.CpuTicker()
	<- ticker.C
	cpusStat2, _ := core.CpuTicker()

    (*sysStat).CpuArray = []core.CpuStat{}
	for i := 0; i < runtime.NumCPU() + 1; i++ {
		cpuStat := core.CpuStat{}

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
	sysStat.Zone.ZoneTicker()
	wg.Done()
}

func (sysStat *SysStat) Paging(t int, wg *sync.WaitGroup) {
	ticker    := time.NewTicker(time.Millisecond * time.Duration(t))
	pageStat  := core.PageStat{}
	pageStat.PageTicker()
	<- ticker.C
	pageStat2 := core.PageStat{}
	pageStat2.PageTicker()

	sysStat.PageIn  = (pageStat2.PageIn  - pageStat.PageIn)  * int64(os.Getpagesize()) * 1
	sysStat.PageOut = (pageStat2.PageOut - pageStat.PageOut) * int64(os.Getpagesize()) * 1

	wg.Done()
}

func (sysStat *SysStat) Disk(t int, wg *sync.WaitGroup) {
    ticker         := time.NewTicker(time.Millisecond * time.Duration(t))
    diskList, err  := core.DiskTicker()
    if err != nil {
        panic(err)
    }
    <- ticker.C
    diskList2, err := core.DiskTicker()
    if err != nil {
        panic(err)
    }

    (*sysStat).DiskList = []core.DiskStat{}
    for index := 0; index < len(diskList); index ++ {
        diskStat := core.DiskStat{}
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
    netList,  _   := core.NetTicker()
    totalStat1    := netList[len(netList) - 1]
    <- ticker.C
    netList2, _   := core.NetTicker()
    totalStat2    := netList2[len(netList2) - 1]

    (*sysStat).NetList = []core.NetStat{}
    for index, netDev := range netDevices {
        netStat := core.NetStat{}
        netStat.Name = netDev
        netStat.Recv = netList2[index].Recv - netList[index].Recv
        netStat.Send = netList2[index].Send - netList[index].Send

        (*sysStat).NetList = append((*sysStat).NetList, netStat)
    }

    totalStat := core.NetStat{}
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
    swapList, _ := core.SwapTicker()
    (*sysStat).SwapList = swapList
    wg.Done()
}

func (sysStat *SysStat) System(t int, wg *sync.WaitGroup) {
    ticker      := time.NewTicker(time.Millisecond * time.Duration(t))
    systemStat  := core.SystemStat{}
    systemStat.SystemTicker()
    <- ticker.C
    systemStat2 := core.SystemStat{}
    systemStat2.SystemTicker()

    sysStat.Interrupt     = systemStat2.Interrupt     - systemStat.Interrupt
    sysStat.ContextSwitch = systemStat2.ContextSwitch - systemStat.ContextSwitch

    wg.Done()
}

func (sysStat *SysStat) AllSocket(wg *sync.WaitGroup) {
    sysStat.Socket.SocketTicker()
    sysStat.RawSocket.RawSocketTicker()
    sysStat.UnixSocket.UnixSocketTicker() 
    sysStat.TCP.TCPTicker()
    sysStat.UDP.UDPTicker() 
    wg.Done()
} 

func (sysStat *SysStat) FileSystem(wg *sync.WaitGroup) {
    sysStat.FileSystemTicker()
    wg.Done()
}

func (sysStat *SysStat) IO(t int, wg *sync.WaitGroup) {
    ticker         := time.NewTicker(time.Millisecond * time.Duration(t))
    ioList, err  := core.IOTicker()
    if err != nil {
        panic(err)
    }
    <- ticker.C
    ioList2, err := core.IOTicker()
    if err != nil {
        panic(err)
    }

    (*sysStat).IOList = []core.IOStat{}
    for index := 0; index < len(ioList); index ++ {
        ioStat := core.IOStat{}
        ioStat.Name  = ioList[index].Name
        ioStat.Read  = (ioList2[index].Read  - ioList[index].Read) * 1.0
        ioStat.Write = (ioList2[index].Write - ioList[index].Write) * 1.0

        (*sysStat).IOList = append((*sysStat).IOList, ioStat)
    }
    wg.Done()
}

func (sysStat *SysStat) AIO_(wg *sync.WaitGroup) {
    sysStat.AIO.AIOTicker()
    wg.Done()
} 

func (sysStat *SysStat) Proc_(t int, wg *sync.WaitGroup) {
    ticker    := time.NewTicker(time.Millisecond * time.Duration(t))
    procStat  := core.ProcStat{}
    procStat.ProcTicker()
    <- ticker.C
    procStat2 := core.ProcStat{}
    procStat2.ProcTicker()

    sysStat.Proc.Running   = procStat2.Running
    sysStat.Proc.Blocked   = procStat2.Blocked
    sysStat.Proc.Processes = procStat2.Processes - procStat.Processes
    wg.Done()
}

func (sysStat *SysStat) IPC_(wg *sync.WaitGroup) {
    sysStat.IPC.IPCTicker()
    wg.Done()
}

func (sysStat *SysStat) LockInfo(wg *sync.WaitGroup) {
    sysStat.Lock.LockTicker()
    wg.Done()
}

func (sysStat *SysStat) VMInfo(t int, wg *sync.WaitGroup) {
    ticker  := time.NewTicker(time.Millisecond * time.Duration(t))
    vmStat  := core.VMStat{}
    vmStat.VMTicker()
    <- ticker.C
    vmStat2 := core.VMStat{}
    vmStat2.VMTicker()

    sysStat.VM.PgMajFault = vmStat2.PgMajFault - vmStat.PgMajFault 
    sysStat.VM.PgFault    = vmStat2.PgFault    - vmStat.PgFault 
    sysStat.VM.PgFree     = vmStat2.PgFree     - vmStat.PgFree 
    sysStat.VM.PgAlloc    = vmStat2.PgAlloc    - vmStat.PgAlloc
    wg.Done()
}

func (sysStat *SysStat) Run(t int) {
    writer        := uilive.New()
    writer.Start()

    // fmt.Printf("| --- datetime ---- | ---  cpu(%%) --- | --------- memory usage ---------- | --- paging ---- | - disk total -- | -- net total -- | ---- load avg ---- | ---- swap ----- | --- system ---- |\n")
    // fmt.Printf("|           datetime| user|  sys| idel|    used|    free| buffers|  cached|      in|     out|      in|     out|    recv|    send| load| load5| load15|    used|    free|    intr|     csw|\n")

    for {
        startT := time.Now()
        var wg sync.WaitGroup
        wg.Add(17)

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
        go sysStat.AllSocket(&wg)
        go sysStat.FileSystem(&wg)
        go sysStat.IO(t, &wg)
        go sysStat.AIO_(&wg)
        go sysStat.Proc_(t, &wg)
        go sysStat.IPC_(&wg)
        go sysStat.LockInfo(&wg)
        go sysStat.VMInfo(t, &wg)
        wg.Wait()

        // diskListLength := len((*sysStat).DiskList)
        // netListLength  := len((*sysStat).NetList)
        // swapListLength := len((*sysStat).SwapList)
        tc := time.Since(startT)
        // fmt.Fprintf(writer, "|%s|%5.2f|%5.2f|%5.2f|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%8s|%5.2f|%6.2f|%7.2f|%8s|%8s|%8s|%8s|\n",
        //     time.Time((*sysStat).DateTime).Format("2006-01-02 15:04:05"),
        //     (*sysStat).CpuArray[0].User, (*sysStat).CpuArray[0].System, (*sysStat).CpuArray[0].Idle,
        //     utils.ByteCountSI(int64((*sysStat).Used)), utils.ByteCountSI(int64((*sysStat).Free)), utils.ByteCountSI(int64((*sysStat).Buffers)), utils.ByteCountSI(int64((*sysStat).Cached)),
        //     utils.ByteCountSI((*sysStat).PageIn), utils.ByteCountSI((*sysStat).PageOut),
        //     utils.ByteCountSI(int64((*sysStat).DiskList[diskListLength - 1].Read)), utils.ByteCountSI(int64((*sysStat).DiskList[diskListLength - 1].Write)),
        //     utils.ByteCountSI(int64((*sysStat).NetList[netListLength - 1].Recv)), utils.ByteCountSI(int64((*sysStat).NetList[netListLength - 1].Send)),
        //     (*sysStat).Load1, (*sysStat).Load5, (*sysStat).Load15,
        //     utils.ByteCountSI(int64((*sysStat).SwapList[swapListLength - 1].Used)), utils.ByteCountSI(int64((*sysStat).SwapList[swapListLength - 1].Free)),
        //     utils.ByteCountSI(int64(sysStat.Interrupt)), utils.ByteCountSI(int64(sysStat.ContextSwitch)))
        // fmt.Fprintf(writer, "---------------------------------------------\n") 
        // fmt.Fprintf(writer, "|%5d|%5d|%5d|%5d|%5d|%5d|\n", (*sysStat).Socket.Total, (*sysStat).Socket.TCP, (*sysStat).Socket.UDP, (*sysStat).Socket.RAW, (*sysStat).Socket.FRAG, (*sysStat).Socket.Other)
        // fmt.Fprintf(writer, "---------------------------------------------\n") 
        // fmt.Fprintf(writer, "|%5d|%5d|\n", (*sysStat).UsingFileHandle, (*sysStat).UsingInode)
        // fmt.Fprintf(writer, "---------------------------------------------\n") 
        // for _, ioDev := range (*sysStat).IOList {
        //     fmt.Fprintf(writer, "|%8s|%5.2f|%5.2f|\n", ioDev.Name, ioDev.Read, ioDev.Write)
        // }
        // fmt.Fprintf(writer, "---------------------------------------------\n") 
        // fmt.Fprintf(writer, "|%8d|\n", (*sysStat).AIO.Requests)
        // fmt.Fprintf(writer, "---------------------------------------------\n")
        // fmt.Fprintf(writer, "|%5f|%5f|%5f|\n", (*sysStat).Proc.Running, (*sysStat).Proc.Blocked, (*sysStat).Proc.Processes)
        // fmt.Fprintf(writer, "---------------------------------------------\n")
        // fmt.Fprintf(writer, "|%8d|%8d|%8d|\n", (*sysStat).IPC.MessageQueue, (*sysStat).IPC.Semaphore, (*sysStat).IPC.SharedMemory)
        // fmt.Fprintf(writer, "---------------------------------------------\n")
        // fmt.Fprintf(writer, "|%8s|%8s|%8s|%8s|\n", utils.ByteCountSI((*sysStat).Zone.DMA32Free), utils.ByteCountSI((*sysStat).Zone.DMA32High), utils.ByteCountSI((*sysStat).Zone.NormalFree), utils.ByteCountSI((*sysStat).Zone.NormalHigh))
        // fmt.Fprintf(writer, "---------------------------------------------\n")
        // fmt.Fprintf(writer, "|%8d|%8d|%8d|%8d|\n", (*sysStat).Lock.Posix, (*sysStat).Lock.Flock, (*sysStat).Lock.Read, (*sysStat).Lock.Write)
        fmt.Fprintf(writer, "---------------------------------------------\n")
        fmt.Fprintf(writer, "|%8d|%8d|%8d|%8d|\n", (*sysStat).VM.PgMajFault, (*sysStat).VM.PgFault, (*sysStat).VM.PgAlloc, (*sysStat).VM.PgFree)
        fmt.Fprintf(writer, "time const = %v\n", tc) 
    }
    writer.Stop()
}

