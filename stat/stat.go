package stat

import (
    "fmt"
    "os"
    "sync"
    "time"
    "path"
    "strconv"
    "strings"
    "runtime"
    "encoding/csv"

    "github.com/gosuri/uilive"

    utils "godstat/utils"
    core  "godstat/core"
)

type SysStat struct {
    DateTime   utils.FormatTime    `json:"datetime"`
    Epoch      int64               `json:"epoch"`
    CpuArray   []core.CpuStat      `json:"cpuList"`
    Memory     core.MemoryStat 
    Page       core.PageStat 
    DiskList   []core.DiskStat     `json:"diskList"`
    NetList    []core.NetStat      `json:"netList"`
    LoadAvg    core.LoadStat
    SwapList   []core.SwapStat     `json:"swapList"`
    System     core.SystemStat
    Socket     core.SocketStat
    RawSocket  core.RawSocketStat
    UnixSocket core.UnixSocketStat 
    TCP        core.TCPStat 
    UDP        core.UDPStat
    FileSystem core.FileSystemStat
    IOList     []core.IOStat       `json:"ioList"`
    AIO        core.AIOStat
    Proc       core.ProcStat
    IPC        core.IPCStat
    Zone       core.ZoneStat 
    Lock       core.LockStat 
    VM         core.VMStat 
}

func (sysStat *SysStat) getCpuUtilization(t int, wg *sync.WaitGroup) {
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

func (sysStat *SysStat) getMemory(wg *sync.WaitGroup) {
	sysStat.Memory.MemoryTicker()
	wg.Done()
}

func (sysStat *SysStat) getZones(wg *sync.WaitGroup) {
	sysStat.Zone.ZoneTicker()
	wg.Done()
}

func (sysStat *SysStat) getPaging(t int, wg *sync.WaitGroup) {
	ticker    := time.NewTicker(time.Millisecond * time.Duration(t))
	pageStat  := core.PageStat{}
	pageStat.PageTicker()
	<- ticker.C
	pageStat2 := core.PageStat{}
	pageStat2.PageTicker()

	sysStat.Page.PageIn  = (pageStat2.PageIn  - pageStat.PageIn)  * int64(os.Getpagesize()) * 1
	sysStat.Page.PageOut = (pageStat2.PageOut - pageStat.PageOut) * int64(os.Getpagesize()) * 1

	wg.Done()
}

func (sysStat *SysStat) getDisk(t int, wg *sync.WaitGroup) {
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
        diskStat      := core.DiskStat{}
        diskStat.Name  = diskList[index].Name
        diskStat.Read  = (diskList2[index].Read  - diskList[index].Read) * 512.0
        diskStat.Write = (diskList2[index].Write - diskList[index].Write) * 512.0

        (*sysStat).DiskList = append((*sysStat).DiskList, diskStat)
    }

    wg.Done()
}

func (sysStat *SysStat) getNet(t int, wg *sync.WaitGroup) {
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

func (sysStat *SysStat) getLoadAvg(wg *sync.WaitGroup) {
    sysStat.LoadAvg.LoadTicker()
    wg.Done()
}

func (sysStat *SysStat) getSwap(wg *sync.WaitGroup) {
    swapList, _ := core.SwapTicker()
    (*sysStat).SwapList = swapList
    wg.Done()
}

func (sysStat *SysStat) getSystem(t int, wg *sync.WaitGroup) {
    ticker      := time.NewTicker(time.Millisecond * time.Duration(t))
    systemStat  := core.SystemStat{}
    systemStat.SystemTicker()
    <- ticker.C
    systemStat2 := core.SystemStat{}
    systemStat2.SystemTicker()

    sysStat.System.Interrupt     = systemStat2.Interrupt     - systemStat.Interrupt
    sysStat.System.ContextSwitch = systemStat2.ContextSwitch - systemStat.ContextSwitch

    wg.Done()
}

func (sysStat *SysStat) getSocket(wg *sync.WaitGroup) {
    sysStat.Socket.SocketTicker()
    wg.Done()
}

func (sysStat *SysStat) getRawSocket(wg *sync.WaitGroup) {
    sysStat.RawSocket.RawSocketTicker()
    wg.Done()
}

func (sysStat *SysStat) getUnixSocket(wg *sync.WaitGroup) {
    sysStat.UnixSocket.UnixSocketTicker()
    wg.Done()
}

func (sysStat *SysStat) getTCPSocket(wg *sync.WaitGroup) {
    sysStat.TCP.TCPTicker()
    wg.Done()
}

func (sysStat *SysStat) getUDPSocket(wg *sync.WaitGroup) {
    sysStat.UDP.UDPTicker()
    wg.Done()
}

func (sysStat *SysStat) getFileSystem(wg *sync.WaitGroup) {
    sysStat.FileSystem.FileSystemTicker()
    wg.Done()
}

func (sysStat *SysStat) getIO(t int, wg *sync.WaitGroup) {
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

func (sysStat *SysStat) getAIO(wg *sync.WaitGroup) {
    sysStat.AIO.AIOTicker()
    wg.Done()
} 

func (sysStat *SysStat) getProc(t int, wg *sync.WaitGroup) {
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

func (sysStat *SysStat) getIPC(wg *sync.WaitGroup) {
    sysStat.IPC.IPCTicker()
    wg.Done()
}

func (sysStat *SysStat) getLock(wg *sync.WaitGroup) {
    sysStat.Lock.LockTicker()
    wg.Done()
}

func (sysStat *SysStat) getVM(t int, wg *sync.WaitGroup) {
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

// 1) out 2) out.csv 3) out.csv.0 4) out.csv.123 5) 123
func getNewFileName(fileName, filePath string) string {
    var newName string 
    tempSlice := strings.Split(fileName, ".")
    extName    := tempSlice[len(tempSlice) - 1]
    digit, err := strconv.ParseInt(extName, 10, 64)
    if err == nil {
        digit  += 1
        newNameSlice := tempSlice[:]
        newNameSlice[len(newNameSlice) - 1] = strconv.FormatInt(digit, 10)
        newName = strings.Join(newNameSlice, ".")
    } else {
        digit   = 0
        newName = fileName + "." + strconv.FormatInt(digit, 10)
    }
    _, err = os.Stat(path.Join(filePath, newName))
    if err == nil {
        return getNewFileName(newName, filePath)
    }
    return newName
}

func (sysStat *SysStat) Run(t int, 
                            cpuSlice,
                            diskSlice,
                            netSlice, 
                            swapSlice []string, 
                            isPage, 
                            isLoad, 
                            isMem, 
                            isProc, 
                            isIO, 
                            isTime, 
                            isEpoch, 
                            isSys, 
                            isFS, 
                            isAIO, 
                            isIPC, 
                            isLock, 
                            isRAW, 
                            isSocket, 
                            isTCP, 
                            isUDP, 
                            isUnix, 
                            isVM, 
                            isZones bool,
                            outCSVPath string) { 
    var csvWriter  *csv.Writer
    var csvTitle   []string 

    diskNames, err := utils.DiskDev()
    if err != nil {
        fmt.Println(err)
    }
    netNames,  err := utils.NetDev()
    if err != nil {
        fmt.Println(err)
    }
    swapNames, err := utils.SwapList()
    if err != nil {
        fmt.Println(err)
    }

    writer    := uilive.New()
    writer.Start()
    
    if outCSVPath != "" {
        // 1) out.csv 2) ./out.csv 3) /tmp/out.csv 4) /tmp/test/out.csv (/tmp/test is not exist.)
        outCSVPathSplit := strings.Split(outCSVPath, "/")
        outCSVFileName  := outCSVPathSplit[len(outCSVPathSplit) - 1]
        outCSVDir       := strings.Replace(outCSVPath, outCSVFileName, "", -1) // "", "./", "/tmp/test"
        
        if outCSVDir != "" {
            _, err := os.Stat(outCSVDir)
            if err != nil {
                if os.IsExist(err) {
                    fmt.Println(os.IsExist(err))
                }
                err = os.MkdirAll(outCSVDir, os.ModePerm)
                if err != nil {
                    fmt.Println(err)
                }
            }
        } else {
            outCSVDir = "./"
        }
        _, err := os.Stat(outCSVPath)
        if err == nil {
            // if old file exist, rename it.
            newFileName := getNewFileName(outCSVFileName, outCSVDir)
            os.Rename(path.Join(outCSVDir, outCSVFileName), path.Join(outCSVDir, newFileName))
            fmt.Printf("rename exist file: %s --> %s.\n", path.Join(outCSVDir, outCSVFileName), 
                                                          path.Join(outCSVDir, newFileName))
        }
        fmt.Printf("out csv file: %s\n", outCSVPath)
        csvF, err := os.Create(outCSVPath)
        if err != nil {
            fmt.Println(err)            
        }
        defer csvF.Close()
        csvWriter = csv.NewWriter(csvF)

        if len(cpuSlice) !=0  {for _, _ = range cpuSlice   {csvTitle = append(csvTitle, []string{"cpu", "user", "system", "idle", "iowait", "steal"}...)}}
        if isMem              {csvTitle = append(csvTitle, []string{"used(memory)", "usedPCT(memory)", "free(memory)", "buffers(memory)", "cached(memory)"}...)}
        if isPage             {csvTitle = append(csvTitle, []string{"pageIn", "pageOut"}...)}
        if len(diskSlice)!=0  {for _, _ = range diskSlice  {csvTitle = append(csvTitle, []string{"diskName(disk)", "read(disk)", "write(disk)"}...)}}
        if len(netSlice) !=0  {for _, _ = range netSlice   {csvTitle = append(csvTitle, []string{"net", "recv", "send"}...)}}
        if isLoad             {csvTitle = append(csvTitle, []string{"load1", "load5", "load15"}...)}
        if len(swapSlice)!=0  {for _, _ = range swapSlice  {csvTitle = append(csvTitle, []string{"swap", "used", "free"}...)}}
        if isSys              {csvTitle = append(csvTitle, []string{"interrupt", "contextSwitch"}...)}
        if isSocket           {csvTitle = append(csvTitle, []string{"total(socket)", "tcp", "udp", "raw", "frag", "other"}...)}
        if isRAW              {csvTitle = append(csvTitle, []string{"rawSocket"}...)}
        if isUnix             {csvTitle = append(csvTitle, []string{"dataGram(unix)", "stream(unix)", "established(unix)", "listen(unix)"}...)}
        if isTCP              {csvTitle = append(csvTitle, []string{"listen(tcp)", "established(tcp)", "syn(tcp)", "timeWait(tcp)", "close(tcp)"}...)}
        if isUDP              {csvTitle = append(csvTitle, []string{"listen(udp)", "established(udp)"}...)}
        if isFS               {csvTitle = append(csvTitle, []string{"usingFileHandle", "usingInode"}...)}
        if isIO               {for _, _ = range diskSlice  {csvTitle = append(csvTitle, []string{"diskName(io)", "read(io)", "write(io)"}...)}}
        if isAIO              {csvTitle = append(csvTitle, []string{"requests"}...)}
        if isProc             {csvTitle = append(csvTitle, []string{"running", "blocked", "processes"}...)}
        if isIPC              {csvTitle = append(csvTitle, []string{"messageQueue", "semaphore", "sharedMemory"}...)}
        if isZones            {csvTitle = append(csvTitle, []string{"dma2Free", "dma32High", "normalFree", "normalHigh"}...)}
        if isLock             {csvTitle = append(csvTitle, []string{"posix", "flock", "read", "write"}...)}
        if isVM               {csvTitle = append(csvTitle, []string{"pgMajFault", "pgFault", "pgAlloc", "pgFree"}...)}
        if isEpoch            {csvTitle = append(csvTitle, []string{"epoch"}...)}

        csvTitle = append(csvTitle, []string{"datetime"}...)
        
        csvWriter.Write(csvTitle)
        csvWriter.Flush()
    }

    for {
        // startT  := time.Now()
        var csvData []string
        var wg sync.WaitGroup
        wg.Add(22)

        go func(sysStat *SysStat, wg *sync.WaitGroup) {
            sysStat.DateTime = utils.FormatTime(time.Now())
            wg.Done()
        }(sysStat, &wg)
        go func(sysStat *SysStat, wg *sync.WaitGroup) {
            sysStat.Epoch = time.Now().Unix()
            wg.Done()
        }(sysStat, &wg)
        go sysStat.getCpuUtilization(t, &wg)
        go sysStat.getDisk(t, &wg)
        go sysStat.getNet(t, &wg)
        go sysStat.getSwap(&wg)
        go sysStat.getPaging(t, &wg)
        go sysStat.getLoadAvg(&wg)
        go sysStat.getMemory(&wg)
        go sysStat.getProc(t, &wg)
        go sysStat.getIO(t, &wg)
        go sysStat.getSystem(t, &wg)
        go sysStat.getFileSystem(&wg)
        go sysStat.getAIO(&wg)
        go sysStat.getIPC(&wg)
        go sysStat.getLock(&wg)
        go sysStat.getRawSocket(&wg)
        go sysStat.getSocket(&wg)
        go sysStat.getTCPSocket(&wg)
        go sysStat.getUDPSocket(&wg)
        go sysStat.getUnixSocket(&wg)
        go sysStat.getVM(t, &wg)
        wg.Wait()

        // tc := time.Since(startT)
         
        // cpu
        if len(cpuSlice) != 0 {
            fmt.Fprintf(writer, "|%16s cpu usage %14s|\n", "", "")
            fmt.Fprintf(writer, "|%6s|%6s|%6s|%6s|%6s|%6s|\n", "cpu", "user", "system", "idle", "iowait", "steal")
            for _, cpuName := range cpuSlice {
                var index int64
                var err error
                if cpuName == "total" {
                    index = 0
                } else {
                    index, err = strconv.ParseInt(cpuName, 0, 64)
                    if err != nil {
                        fmt.Fprintf(writer, "%s", err.Error())
                    }
                    index += 1
                }
                cpuStat := (*sysStat).CpuArray[index]
                fmt.Fprintf(writer, "|%6s|%6.2f|%6.2f|%6.2f|%6.2f|%6.2f|\n", cpuStat.CPU, cpuStat.User, cpuStat.System, cpuStat.Idle, cpuStat.Iowait, cpuStat.Steal)
                
                if outCSVPath != "" {
                    //csvTitle = append(csvTitle, []string{"cpu", "user", "system", "idle", "iowait", "steal"}...)
                    csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%f,%f,%f,%f,%f", cpuStat.CPU, cpuStat.User, cpuStat.System, cpuStat.Idle, cpuStat.Iowait, cpuStat.Steal), ",")...)
                }
            } 
        }

        // memory
        if isMem {
            fmt.Fprintf(writer, "|%17s memory %16s|\n", "", "")
            fmt.Fprintf(writer, "|%8s|%8s|%7s|%7s|%7s|\n", "used", "usedPCT", "free", "buffers", "cached")
            fmt.Fprintf(writer, "|%8s|%7.2f%%|%7s|%7s|%7s|\n", utils.ByteCountSI(int64((*sysStat).Memory.Used)), (*sysStat).Memory.UsedPercent, utils.ByteCountSI(int64((*sysStat).Memory.Free)), utils.ByteCountSI(int64((*sysStat).Memory.Buffers)), utils.ByteCountSI(int64((*sysStat).Memory.Cached)))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%f,%s,%s,%s", utils.ByteCountSI(int64((*sysStat).Memory.Used)), (*sysStat).Memory.UsedPercent, utils.ByteCountSI(int64((*sysStat).Memory.Free)), utils.ByteCountSI(int64((*sysStat).Memory.Buffers)), utils.ByteCountSI(int64((*sysStat).Memory.Cached))), ",")...)
            }
        }

        // page
        if isPage {
            fmt.Fprintf(writer, "|%18s page %17s|\n", "", "")
            fmt.Fprintf(writer, "|%20s|%20s|\n", "pageIn", "pageOut")
            fmt.Fprintf(writer, "|%20d|%20d|\n", (*sysStat).Page.PageIn, (*sysStat).Page.PageOut)
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d", (*sysStat).Page.PageIn, (*sysStat).Page.PageOut), ",")...)
            }
        }

        // disk
        if len(diskSlice) != 0 {
            fmt.Fprintf(writer, "|%18s disk %17s|\n", "", "")
            fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", "disk", "read", "write")
            for _, diskName := range diskSlice {
                var index int64 
                if diskName == "total" {
                    index = int64(len((*sysStat).DiskList) - 1)
                } else {
                    index = int64(utils.StringsContains(diskNames, diskName))
                }
                diskStat := (*sysStat).DiskList[index]
                fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", diskStat.Name, utils.ByteCountSI(int64(diskStat.Read)), utils.ByteCountSI(int64(diskStat.Write)))
                if outCSVPath != "" {
                    csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%s,%s", diskStat.Name, utils.ByteCountSI(int64(diskStat.Read)), utils.ByteCountSI(int64(diskStat.Write))), ",")...)
                }
            }
        }
       
        // net
        if len(netSlice) != 0 {
            fmt.Fprintf(writer, "|%18s net %18s|\n", "", "")
            fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", "net", "recv", "send")
            for _, netName := range netSlice {
                var index int64 
                if netName == "total" {
                    index = int64(len((*sysStat).NetList) - 1)
                } else {
                    index = int64(utils.StringsContains(netNames, netName))
                }
                netStat := (*sysStat).NetList[index]
                fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", netStat.Name, utils.ByteCountSI(int64(netStat.Recv)), utils.ByteCountSI(int64(netStat.Send)))
                if outCSVPath != "" {
                    csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%s,%s", netStat.Name, utils.ByteCountSI(int64(netStat.Recv)), utils.ByteCountSI(int64(netStat.Send))), ",")...)
                }
            }
        }
       
        // loadavg 
        if isLoad {
            fmt.Fprintf(writer, "|%16s loadAvg %16s|\n", "", "")
            fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", "load1", "load5", "load15")
            fmt.Fprintf(writer, "|%13.2f|%13.2f|%13.2f|\n", (*sysStat).LoadAvg.Load1, (*sysStat).LoadAvg.Load5, (*sysStat).LoadAvg.Load15)
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%f,%f,%f", (*sysStat).LoadAvg.Load1, (*sysStat).LoadAvg.Load5, (*sysStat).LoadAvg.Load15), ",")...)
            }
        }

        // swap 
        if len(swapSlice) != 0 {
            fmt.Fprintf(writer, "|%18s swap %17s|\n", "", "")
            fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", "swap", "used", "free")
            for _, swapName := range swapSlice {
                var index int64 
                if swapName == "total" {
                    index = int64(len((*sysStat).SwapList) - 1)
                } else {
                    index = int64(utils.StringsContains(swapNames, swapName))
                }
                swapStat := (*sysStat).SwapList[index]
                fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", swapStat.Name, utils.ByteCountSI(int64(swapStat.Used)), utils.ByteCountSI(int64(swapStat.Free)))
                if outCSVPath != "" {
                    csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%s,%s", swapStat.Name, utils.ByteCountSI(int64(swapStat.Used)), utils.ByteCountSI(int64(swapStat.Free))), ",")...)
                }
            }
        }
        // system 
        if isSys {
            fmt.Fprintf(writer, "|%17s system %16s|\n", "", "")
            fmt.Fprintf(writer, "|%20s|%20s|\n", "interrupt", "contextSwitch")
            fmt.Fprintf(writer, "|%20d|%20d|\n", int64((*sysStat).System.Interrupt), int64((*sysStat).System.ContextSwitch))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d", int64((*sysStat).System.Interrupt), int64((*sysStat).System.ContextSwitch)), ",")...)
            }
        }
        
        // socket 
        if isSocket {
            fmt.Fprintf(writer, "|%17s socket %16s|\n", "", "")
            fmt.Fprintf(writer, "|%6s|%6s|%6s|%6s|%6s|%6s|\n", "total", "tcp", "udp", "raw", "frag", "other")
            fmt.Fprintf(writer, "|%6d|%6d|%6d|%6d|%6d|%6d|\n", 
                                int64((*sysStat).Socket.Total), 
                                int64((*sysStat).Socket.TCP), 
                                int64((*sysStat).Socket.UDP), 
                                int64((*sysStat).Socket.RAW), 
                                int64((*sysStat).Socket.FRAG),
                                int64((*sysStat).Socket.Other))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d,%d,%d,%d,%d",
                                int64((*sysStat).Socket.Total), 
                                int64((*sysStat).Socket.TCP), 
                                int64((*sysStat).Socket.UDP), 
                                int64((*sysStat).Socket.RAW), 
                                int64((*sysStat).Socket.FRAG),
                                int64((*sysStat).Socket.Other)), ",")...)
            }
        }

        // raw socket 
        if isRAW {
            fmt.Fprintf(writer, "|%15s raw socket %14s|\n", "", "")
            fmt.Fprintf(writer, "|%41s|\n", "rawSocket")
            fmt.Fprintf(writer, "|%41d|\n", int64((*sysStat).RawSocket.NumSockets))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d",int64((*sysStat).RawSocket.NumSockets)), ",")...)
            }
        }

        // unix socket
        if isUnix {
            fmt.Fprintf(writer, "|%14s unix socket %14s|\n", "", "")
            fmt.Fprintf(writer, "|%9s|%9s|%11s|%9s|\n", "dataGram", "stream", "established", "listen")
            fmt.Fprintf(writer, "|%9d|%9d|%11d|%9d|\n", 
                                int64((*sysStat).UnixSocket.DataGram),
                                int64((*sysStat).UnixSocket.Stream),
                                int64((*sysStat).UnixSocket.Established),
                                int64((*sysStat).UnixSocket.Listen))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d,%d,%d",
                                int64((*sysStat).UnixSocket.DataGram),
                                int64((*sysStat).UnixSocket.Stream),
                                int64((*sysStat).UnixSocket.Established),
                                int64((*sysStat).UnixSocket.Listen)), ",")...)
            }
        }

        // tcp 
        if isTCP {
            fmt.Fprintf(writer, "|%15s tcp socket %14s|\n", "", "")
            fmt.Fprintf(writer, "|%6s|%11s|%6s|%8s|%6s|\n", "listen", "established", "syn", "timeWait", "close")
            fmt.Fprintf(writer, "|%6d|%11d|%6d|%8d|%6d|\n", 
                                int64((*sysStat).TCP.Listen),
                                int64((*sysStat).TCP.Established),
                                int64((*sysStat).TCP.SynSent) + int64((*sysStat).TCP.SynRecv) + int64((*sysStat).TCP.LastAck),
                                int64((*sysStat).TCP.TimeWait),
                                int64((*sysStat).TCP.FinWait1) + int64((*sysStat).TCP.FinWait2) + int64((*sysStat).TCP.Close) + int64((*sysStat).TCP.CloseWait) + int64((*sysStat).TCP.Closing))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d,%d,%d,%d",
                                int64((*sysStat).TCP.Listen),
                                int64((*sysStat).TCP.Established),
                                int64((*sysStat).TCP.SynSent) + int64((*sysStat).TCP.SynRecv) + int64((*sysStat).TCP.LastAck),
                                int64((*sysStat).TCP.TimeWait),
                                int64((*sysStat).TCP.FinWait1) + int64((*sysStat).TCP.FinWait2) + int64((*sysStat).TCP.Close) + int64((*sysStat).TCP.CloseWait) + int64((*sysStat).TCP.Closing)), ",")...)
            }
        }

        // udp
        if isUDP {
            fmt.Fprintf(writer, "|%15s udp socket %14s|\n", "", "")
            fmt.Fprintf(writer, "|%20s|%20s|\n", "listen", "established")
            fmt.Fprintf(writer, "|%20d|%20d|\n",
                                int64((*sysStat).UDP.Listen),
                                int64((*sysStat).UDP.Established))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d",
                                int64((*sysStat).UDP.Listen),
                                int64((*sysStat).UDP.Established)), ",")...)
            }
        }

        // filesystem
        if isFS {
            fmt.Fprintf(writer, "|%15s filesystem %14s|\n", "", "")
            fmt.Fprintf(writer, "|%20s|%20s|\n", "usingFileHandle", "usingInode")
            fmt.Fprintf(writer, "|%20d|%20d|\n",
                                int64((*sysStat).FileSystem.UsingFileHandle),
                                int64((*sysStat).FileSystem.UsingInode))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d",
                                int64((*sysStat).FileSystem.UsingFileHandle),
                                int64((*sysStat).FileSystem.UsingInode)), ",")...)
            }
        }

        // io 
        if isIO {
            fmt.Fprintf(writer, "|%19s io %18s|\n", "", "")
            fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", "io", "read", "write")
            for _, diskName := range diskSlice {
                var index int64 
                if diskName == "total" {
                    index = int64(len((*sysStat).IOList) - 1)
                } else {
                    index = int64(utils.StringsContains(diskNames, diskName))
                }
                ioStat := (*sysStat).IOList[index]
                fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", 
                                    ioStat.Name, 
                                    utils.ByteCountSI(int64(ioStat.Read)), 
                                    utils.ByteCountSI(int64(ioStat.Write)))
                if outCSVPath != "" {
                    csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%s,%s",
                                    ioStat.Name, 
                                    utils.ByteCountSI(int64(ioStat.Read)), 
                                    utils.ByteCountSI(int64(ioStat.Write))), ",")...)
                 }
            }
        }

        // aio
        if isAIO {
            fmt.Fprintf(writer, "|%18s aio %18s|\n", "", "")
            fmt.Fprintf(writer, "|%41s|\n", "requests")
            fmt.Fprintf(writer, "|%41d|\n", int64((*sysStat).AIO.Requests))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d",
                                int64((*sysStat).AIO.Requests)), ",")...)
             }
        }


        // proc
        if isProc {
            fmt.Fprintf(writer, "|%18s proc %17s|\n", "", "")
            fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", "running", "blocked", "processes")
            fmt.Fprintf(writer, "|%13d|%13d|%13d|\n", 
                        int64((*sysStat).Proc.Running),
                        int64((*sysStat).Proc.Blocked),
                        int64((*sysStat).Proc.Processes))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d,%d",
                        int64((*sysStat).Proc.Running),
                        int64((*sysStat).Proc.Blocked),
                        int64((*sysStat).Proc.Processes)), ",")...)
             }
        }

        // ipc
        if isIPC {
            fmt.Fprintf(writer, "|%18s proc %17s|\n", "", "")
            fmt.Fprintf(writer, "|%13s|%13s|%13s|\n", "messageQueue", "semaphore", "sharedMemory")
            fmt.Fprintf(writer, "|%13d|%13d|%13d|\n", 
                        int64((*sysStat).IPC.MessageQueue),
                        int64((*sysStat).IPC.Semaphore),
                        int64((*sysStat).IPC.SharedMemory))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d,%d",
                        int64((*sysStat).IPC.MessageQueue),
                        int64((*sysStat).IPC.Semaphore),
                        int64((*sysStat).IPC.SharedMemory)), ",")...)
             }
        }

        // zones
        if isZones {
            fmt.Fprintf(writer, "|%17s zones %17s|\n", "", "")
            fmt.Fprintf(writer, "|%9s|%9s|%10s|%10s|\n", "dma2Free", "dma32High", "normalFree", "normalHigh")
            fmt.Fprintf(writer, "|%9s|%9s|%10s|%10s|\n", 
                        utils.ByteCountSI(int64((*sysStat).Zone.DMA32Free)),
                        utils.ByteCountSI(int64((*sysStat).Zone.DMA32High)),
                        utils.ByteCountSI(int64((*sysStat).Zone.NormalFree)),
                        utils.ByteCountSI(int64((*sysStat).Zone.NormalHigh)))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%s,%s,%s",
                        utils.ByteCountSI(int64((*sysStat).Zone.DMA32Free)),
                        utils.ByteCountSI(int64((*sysStat).Zone.DMA32High)),
                        utils.ByteCountSI(int64((*sysStat).Zone.NormalFree)),
                        utils.ByteCountSI(int64((*sysStat).Zone.NormalHigh))), ",")...)
             }
        }

        // lock
        if isLock {
            fmt.Fprintf(writer, "|%18s lock %17s|\n", "", "")
            fmt.Fprintf(writer, "|%10s|%10s|%9s|%9s|\n", "posix", "flock", "read", "write")
            fmt.Fprintf(writer, "|%10s|%10s|%9s|%9s|\n", 
                        utils.ByteCountSI(int64((*sysStat).Lock.Posix)),
                        utils.ByteCountSI(int64((*sysStat).Lock.Flock)),
                        utils.ByteCountSI(int64((*sysStat).Lock.Read)),
                        utils.ByteCountSI(int64((*sysStat).Lock.Write)))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%s,%s,%s,%s",
                        utils.ByteCountSI(int64((*sysStat).Lock.Posix)),
                        utils.ByteCountSI(int64((*sysStat).Lock.Flock)),
                        utils.ByteCountSI(int64((*sysStat).Lock.Read)),
                        utils.ByteCountSI(int64((*sysStat).Lock.Write))), ",")...)
             }
        }

        // vm
        if isVM {
            fmt.Fprintf(writer, "|%19s vm %18s|\n", "", "")
            fmt.Fprintf(writer, "|%10s|%10s|%9s|%9s|\n", "pgMajFault", "pgFault", "pgAlloc", "pgFree")
            fmt.Fprintf(writer, "|%10d|%10d|%9d|%9d|\n", 
                        int64((*sysStat).VM.PgMajFault),
                        int64((*sysStat).VM.PgFault),
                        int64((*sysStat).VM.PgAlloc),
                        int64((*sysStat).VM.PgFree))
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d,%d,%d,%d",
                        int64((*sysStat).VM.PgMajFault),
                        int64((*sysStat).VM.PgFault),
                        int64((*sysStat).VM.PgAlloc),
                        int64((*sysStat).VM.PgFree)), ",")...)
             }
        }

        if isEpoch {
            fmt.Fprintf(writer, "|Epoch: %34d|\n", (*sysStat).Epoch)
            if outCSVPath != "" {
                csvData  = append(csvData, strings.Split(fmt.Sprintf("%d", (*sysStat).Epoch), ",")...)
             }
        }
        
        if isTime {
            fmt.Fprintf(writer, "|DateTime: %31s|\n", time.Time((*sysStat).DateTime).Format("2006-01-02 15:04:05"))
        }
        if outCSVPath != "" {
           csvData  = append(csvData, strings.Split(fmt.Sprintf("%s", time.Time((*sysStat).DateTime).Format("2006-01-02 15:04:05")), ",")...)
        }
        
        // fmt.Fprintf(writer, "time const = %v\n", tc) 
        if outCSVPath != "" {
            _ = csvWriter.Write(csvData)
            csvWriter.Flush()
        }
    }
    writer.Stop()
}

