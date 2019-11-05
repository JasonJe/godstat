package main

import (
    "os"
    "fmt"
    "time"

    flag "github.com/spf13/pflag"

	stat "godstat/stat"
    info "godstat/info"
)

func main() {
    flag.Usage = func() {
        //flag.SortFlags = false
        fmt.Fprintf(os.Stdout, "Usage of godstat: \n")
        flag.PrintDefaults()
    }
    delayPtr := flag.Int("delay", 1, "time delay.")

    cpuSlicePtr  := flag.StringSliceP("cpuarray",  "C", []string{"total"}, "example: 0,3,total")
    diskSlicePtr := flag.StringSliceP("diskarray", "D", []string{"total"}, "example: total,hda")
    netSlicePtr  := flag.StringSliceP("netarray",  "N", []string{"total"}, "example: eth1,total")
    swapSlicePtr := flag.StringSliceP("swaparray", "S", []string{"total"}, "example: swap1,total")
    //interPtr := flag.StringSliceP("int", "i", []string{""})

    helpPtr   := flag.BoolP("help",       "h", false, "help")
    infoPtr   := flag.BoolP("info",       "i", false, "show system information.")

    outCSVPtr := flag.StringP("out", "o", "", "write CSV output to file. example: --out=./out.csv") 
    flag.Lookup("out").NoOptDefVal = fmt.Sprintf("%s", time.Now().Format("2006-01-02")) + ".csv"

    cpuPtr    := flag.BoolP("cpu",        "c", false, "enable cpu stats.")
    diskPtr   := flag.BoolP("disk",       "d", false, "enable disk stats.")
    netPtr    := flag.BoolP("net",        "n", false, "enable net stats.")
    swapPtr   := flag.BoolP("swap",       "s", false, "enable swap stats.")
    pagePtr   := flag.BoolP("page",       "g", false, "enable page stats.")
    loadPtr   := flag.BoolP("load",       "l", false, "enable load stats.")
    memPtr    := flag.BoolP("mem",        "m", false, "enable memory stats.")
    procPtr   := flag.BoolP("proc",       "p", false, "enable process stats.")
    ioPtr     := flag.BoolP("io",         "r", false, "enable io stats. \n\t(I/O requests completed)")
    sysPtr    := flag.BoolP("sys",        "y", false, "enable system stats.")
    fsPtr     := flag.BoolP("filesystem", "f", false, "enable filesystem stats.")
    timePtr   := flag.BoolP("time",       "t", false, "enable time/date output.")
    epochPtr  := flag.BoolP("epoch",      "T", false, "enable time counter. (Seconds since epoch)")
    
    aioPtr    := flag.Bool("aio",    false, "enable aio stats.")
    ipcPtr    := flag.Bool("ipc",    false, "enable ipc stats.")
    lockPtr   := flag.Bool("lock",   false, "enable lock stats.")
    rawPtr    := flag.Bool("raw",    false, "enable raw  stats.")
    socketPtr := flag.Bool("socket", false, "enable socket stats.")
    tcpPtr    := flag.Bool("tcp",    false, "enable tcp stats.")
    udpPtr    := flag.Bool("udp",    false, "enable udp stats.")
    unixPtr   := flag.Bool("unix",   false, "enable unix stats.")
    vmPtr     := flag.Bool("vm",     false, "enable vm stats.")
    zonesPtr  := flag.Bool("zones",  false, "enable zoneinfo stats.")
     
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
    
    // go run dstat.go -c -C 0,total -m -d -D total,sda -n -N total,enp6s0 -s -S total -y --socket --raw --unix --tcp --udp --filesystem --io --aio --proc --ipc --zones --lock --vm -T -o=./out/out.csv
    sysStat := new(stat.SysStat)
    if *cpuPtr    ||  
       *diskPtr   || 
       *netPtr    || 
       *swapPtr   || 
       *pagePtr   || 
       *loadPtr   || 
       *memPtr    || 
       *procPtr   || 
       *ioPtr     || 
       *sysPtr    || 
       *fsPtr     || 
       *timePtr   || 
       *epochPtr  || 
       *aioPtr    || 
       *ipcPtr    || 
       *lockPtr   || 
       *rawPtr    || 
       *socketPtr || 
       *tcpPtr    || 
       *udpPtr    || 
       *unixPtr   || 
       *vmPtr     || 
       *zonesPtr  {
        if !*cpuPtr {
            *cpuSlicePtr  = []string{}
        }
        if !*diskPtr {
            *diskSlicePtr = []string{} 
        }
        if !*netPtr {
            *netSlicePtr  = []string{}
        }
        if !*swapPtr {
            *swapSlicePtr = []string{}
        }
        sysStat.Run(*delayPtr * 1000,
                    *cpuSlicePtr,
                    *diskSlicePtr,
                    *netSlicePtr,
                    *swapSlicePtr,
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
                    *zonesPtr,
                    *outCSVPtr)
   } else {
        sysStat.Run(1 * 1000,
                    []string{"total"},
                    []string{"total"},
                    []string{"total"},
                    []string{"total"},
                    true,
                    *loadPtr,
                    true,
                    *procPtr,
                    *ioPtr,
                    true,
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
                    *zonesPtr,
                    *outCSVPtr)
   }
}
