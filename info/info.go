package main

import (
    "fmt"
    "os"
    //"bytes"
    //"encoding/binary"
    "path" 
    "unsafe"
    "syscall"
    "strings"
    "runtime"
    "io/ioutil"

    utils "../utils"
)

type SystemConfig interface {
    GetConfig() error
}

type KernelConfig struct {
    Release       string `json:"kernelRelease"`
    Version       string `json:"kernelVersion"`
}

func (kernelConfig *KernelConfig) GetConfig() error {
    release, _ := utils.ReadLines("/proc/sys/kernel/osrelease")
    version, _ := utils.ReadLines("/proc/sys/kernel/version")
    
    var uname syscall.Utsname
    if err := syscall.Uname(&uname); err != nil {
        return err
    }

    kernelConfig.Release = release[0]
    kernelConfig.Version = version[0]
    return nil
}

type OsConfig struct {
    Name          string `json:"OSName"`
    Vendor        string `json:"OSVendor"`
    Version       string `json:"OSVersion"`
    Release       string `json:"OSRelease"`
    Architercture string `json:"OSArchitercture"`
}

func (osConfig *OsConfig) GetConfig() error {
    var uname syscall.Utsname
    if err := syscall.Uname(&uname); err != nil {
        return err
    }
    // unsafe.Pointer() 包含任意类型的地址
    // (*[65]byte)(unsafe.Pointer(&uname.Machine)) 将该地址装维 byte 数组
    // strings.TrimRight() 删除字符串右边指定字符
    osConfig.Architercture = strings.TrimRight(string((*[65]byte)(unsafe.Pointer(&uname.Machine))[:]), "\000")
    
    osReleaseFile := "/etc/os-release"
    lines, err := utils.ReadLines(osReleaseFile)
    if err != nil {
        if os.IsNotExist(err) {
            if _, err := os.Stat("/etc/redhat-release"); !os.IsNotExist(err) {
                osReleaseFile = "/etc/redhat-release"
                lines, err = utils.ReadLines(osReleaseFile)
            } else if _, err := os.Stat("/etc/centos-release"); !os.IsNotExist(err) {
                osReleaseFile = "/etc/centos-release"
                lines, err = utils.ReadLines(osReleaseFile)
            } else {
                return err 
            }
        } else {
            return err
        }
    }
    
    for _, line := range lines {
        fields := strings.Split(line, "=")
        switch fields[0] {
        case "NAME": 
            osConfig.Name = fields[1]
        case "ID":
            osConfig.Vendor = fields[1]
        case "VERSION_ID":
            osConfig.Version = fields[1]
        case "PRETTY_NAME":
            osConfig.Release = fields[1]
        }
    }
    return nil
}

type CpuConfig struct {
    Vendor  string `json:"CPUVendor"`
    Model   string `json:"CPUModel"`
    Speed   int    `json:"CPUSpeed"`
    Cache   string `json:"CPUCache"`
    Cpus    int    `json:"CPUs"`
    Cores   int    `json:"CPUCores"`
    Threads int    `json:"CPUThread"`
}

func (cpuConfig *CpuConfig) GetConfig() error {
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
    return nil
}

// type MemoryConfig struct {
//     Type  string `json:"memoryType"`
//     Speed int    `json:"memorySpeed"`
//     Size  int    `json:"memorySize"`
// }
// 
// func (memoryConfig *MemoryConfig) GetConfig() error {
//     mem, err := ioutil.ReadFile("/sys/firmware/dmi/tables/DMI")
//     if err != nil {
//         fmt.Println(err)
//         return err 
//     }
//     fmt.Println(mem)
//     for p := 0; p < len(mem) - 1; {
//         recType := mem[p]
//         recLen := mem[p+1]
// 
//         switch recType {
//         case 17:
//             index := p+0x0c
//             size  := binary.LittleEndian.Uint16(mem[index: index+2])
// 
//             fmt.Println(size)
//         }
//         for p += int(recLen); p < len(mem) - 1; {
//             if bytes.Equal(mem[p:p+2], []byte{0, 0}) {
//                 p += 2
//                 break                                                            
//             }
//             p++
//         }
//     }
//     return nil 
// }

func main() {
    var config SystemConfig 
    
    config = &KernelConfig{}
    config.GetConfig()
    fmt.Println(config)

    config = &OsConfig{}
    err := config.GetConfig()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(config)

    config = &CpuConfig{}
    config.GetConfig()
    fmt.Println(config)

//    config = &MemoryConfig{}
//    config.GetConfig()

    devices, err := ioutil.ReadDir("/sys/class/net") // 读取下面的所有目录、文件
    if err != nil {
        panic(err)
    }
    for _, link := range devices {
        fullpath := path.Join("/sys/class/net", link.Name())
        dev, err := os.Readlink(fullpath)
        if err != nil {
            continue 
        }

        if strings.HasPrefix(dev, "../../devices/virtual/") {
            continue 
        }
    
        config = &NICConfig{Name: link.Name()}
        config.GetConfig()
        fmt.Println(config)

    }
}
