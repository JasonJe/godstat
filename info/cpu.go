package info 

import (
    "fmt"
    "bytes"
    "encoding/binary"
    "errors"
    "strings"
    "runtime"

    "github.com/digitalocean/go-smbios/smbios"
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
    if s, ok := args[0].(*smbios.Structure); ok {
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

        var speedU uint16 
        binary.Read(bytes.NewBuffer(s.Formatted[0x12: 0x14][0:2]), binary.LittleEndian, &speedU)
        cpuConfig.Speed = int(speedU)
        return nil 
    }
    return errors.New("unkown")
}
