package info 

import (
    "fmt"
    "bytes"
    "syscall"
    "strings"
    "runtime"
    "encoding/binary"

    utils "godstat/utils"
)

type CpuConfig struct {
    Vendor  string `json:"cpuVendor"`
    Model   string `json:"cpuModel"`
    Speed   int    `json:"cpuSpeed"`
    Cache   string `json:"cpuCache"`
    Cpus    int    `json:"cpus"`
    Cores   int    `json:"cpuCores"`
    Threads int    `json:"cpuThread"`
}

func (cpuConfig *CpuConfig) GetConfig(args ...interface{}) error {
    var cpuID string 
    lines, _ := utils.ReadLines("/proc/cpuinfo")
    cpu      := make(map[string]bool)
    core     := make(map[string]bool)

    for _, line := range lines {            
        fields := strings.Split(line, ":")
        if len(fields) < 2 {
            continue
        }
        
        switch strings.TrimSpace(fields[0]) {
        case "physical id":
            cpuID            = strings.TrimSpace(fields[1])
            cpu[cpuID]       = true 
        case "core id":
            coreID          := fmt.Sprintf("%s:%s", cpuID, strings.TrimSpace(fields[1]))
            core[coreID]     = true
        case "vendor_id":
            cpuConfig.Vendor = strings.TrimSpace(fields[1])
        case "model name":
            cpuConfig.Model  = strings.TrimSpace(fields[1])
        case "cache size":
            cpuConfig.Cache  = strings.TrimSpace(fields[1])
        }
    }
    cpuConfig.Cpus    = len(cpu)
    cpuConfig.Cores   = len(core) 
    cpuConfig.Threads = runtime.NumCPU()
    
    mem, err := utils.StructureTable()
    if err != nil {
        return err 
    } 
    defer syscall.Munmap(mem) // mmap 将一个文件或者其它对象映射进内存, munmap 解除内存映射

    for p := 0; p < len(mem) - 1; {
        recType := mem[p]
        recLen  := mem[p + 1]
        if recType == 4 {
            speed := binary.LittleEndian.Uint16(mem[p + 0x16: p + 0x16 + 2])
            cpuConfig.Speed = int(speed)
        } else if recType == 127 {
            break
        }
        for p += int(recLen); p < len(mem)-1; {
            if bytes.Equal(mem[p: p + 2], []byte{0, 0}) {
                p += 2
                break
            }
            p++ 
        }
    }
    return nil 
}
