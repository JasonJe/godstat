package main

import (
    "os"
    "fmt"
    
    flag "github.com/spf13/pflag"

	stat "godstat/stat"
    info "godstat/info"
)

func main() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stdout, "Usage of godstat: \n")
        flag.PrintDefaults()
    }
    delayPtr := flag.Int("delay", 1, "Time delay.")

    cpuPtr   := flag.StringSliceP("cpu",  "c", []string{"total"}, "Enable cpu  stats. Example: 0,3,total")
    diskPtr  := flag.StringSliceP("disk", "d", []string{"total"}, "Enable disk stats. Example: total,hda")
    netPtr   := flag.StringSliceP("net",  "n", []string{"total"}, "Enable network stats. Example: eth1,total")
    swapPtr  := flag.StringSliceP("swap", "s", []string{"total"}, "Enable swap stats. Example: swap1,total")
    //interPtr := flag.StringSliceP("int", "i", []string{""})

    helpPtr   := flag.BoolP("help",       "h", false, "help")
    infoPtr   := flag.BoolP("info",       "i", false, "Show system information.")
    pagePtr   := flag.BoolP("page",       "g", true,  "Enable page stats.")
    loadPtr   := flag.BoolP("load",       "l", false, "Enable load stats.")
    memPtr    := flag.BoolP("mem",        "m", true,  "Enable memory stats.")
    procPtr   := flag.BoolP("proc",       "p", false, "Enable process stats.")
    ioPtr     := flag.BoolP("io",         "r", false, "Enable io stats. (I/O requests completed)")
    timePtr   := flag.BoolP("time",       "t", true,  "Enable time/date output.")
    epochPtr  := flag.BoolP("epoch",      "T", false, "Enable time counter. (Seconds since epoch)")
    sysPtr    := flag.BoolP("sys",        "y", true,  "Enable system stats.")
    fsPtr     := flag.BoolP("filesystem", "f", false, "Enable filesystem stats.")
    
    aioPtr    := flag.Bool("aio",    false, "Enable aio stats.")
    ipcPtr    := flag.Bool("ipc",    false, "Enable ipc stats.")
    lockPtr   := flag.Bool("lock",   false, "Enable lock stats.")
    rawPtr    := flag.Bool("raw",    false, "Enable raw stats.")
    socketPtr := flag.Bool("socket", false, "Enable socket stats.")
    tcpPtr    := flag.Bool("tcp",    false, "Enable tcp stats.")
    udpPtr    := flag.Bool("udp",    false, "Enable udp stats.")
    unixPtr   := flag.Bool("unix",   false, "Enable unix stats.")
    vmPtr     := flag.Bool("vm",     false, "Enable vm stats.")
    zonesPtr  := flag.Bool("zones",  false, "Enable zoneinfo stats.")
    
    flag.Parse()
    
    if *helpPtr {
        flag.Usage()
        return 
    }
    
    if *infoPtr {
        for _, info := range info.GetInfoFmt() {
            fmt.Printf(info)
        }
        return 
    }
    // go run dstat.go --aio -c total -d total -t -f -r --ipc -l --lock -m -n total -g -p --raw --socket -s total -y --tcp -t --udp --unix --vm -- zones -T
    sysStat := new(stat.SysStat)    
    sysStat.Run(*delayPtr * 1000,
                *cpuPtr,
                *diskPtr,
                *netPtr,
                *swapPtr,
                *pagePtr,
                *loadPtr,
                *memPtr,
                *procPtr,
                *ioPtr,
                *timePtr,
                *epochPtr,
                *sysPtr,
                *fsPtr,
                *aioPtr,
                *ipcPtr,
                *lockPtr,
                *rawPtr,
                *socketPtr,
                *tcpPtr,
                *udpPtr,
                *unixPtr,
                *vmPtr,
                *zonesPtr)
}
