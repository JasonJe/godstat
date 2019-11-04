package info

import (
    "fmt"
    "os"
    "path" 
    "strings"
    "io/ioutil"

    utils "godstat/utils"
)

type SystemConfig interface {
    GetConfig(args ...interface{}) error
    GetInfoFmt() string 
}

func GetInfoFmt() []string {
    var err error 
    var fmtInfo []string 
    var index int
    var systemConfig SystemConfig
   
    // os
    systemConfig = &OsConfig{}
    err = systemConfig.GetConfig()
    if err != nil {
        fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get os info, %s\n", err))
    }
    osInfo := *systemConfig.(*OsConfig)
    osInfoFmt := osInfo.GetInfoFmt()
    fmtInfo = append(fmtInfo, osInfoFmt)

    // kernel
    systemConfig = &KernelConfig{}
    err = systemConfig.GetConfig()
    if err != nil {
        fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get kernel info, %s\n", err))
    }
    kernelInfo := *systemConfig.(*KernelConfig)
    kernelInfoFmt := kernelInfo.GetInfoFmt()
    fmtInfo = append(fmtInfo, kernelInfoFmt)
     
    // cpu
    systemConfig  = &CpuConfig{} 
    err = systemConfig.GetConfig()
    if err != nil {
        fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get cpu info, %s\n", err))
    }
    cpuInfo    := *systemConfig.(*CpuConfig)
    cpuInfoFmt := cpuInfo.GetInfoFmt()
    fmtInfo = append(fmtInfo, cpuInfoFmt)
    
  
    // memory
    systemConfig = &MemoryConfig{}
    err = systemConfig.GetConfig()
    if err != nil {
        fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get memory info, %s\n", err))
    }
    memoryInfo := *systemConfig.(*MemoryConfig)
    memoryInfoFmt := memoryInfo.GetInfoFmt()
    fmtInfo = append(fmtInfo, memoryInfoFmt)
   
    // nic
    fmtInfo = append(fmtInfo, "NIC Info:\n")
    nicDevs, err := ioutil.ReadDir("/sys/class/net") // 读取下面的所有目录、文件
    if err != nil {
        fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get nic info, %s\n", err))
    }
    index = 1
    for _, link := range nicDevs {
        fullpath := path.Join("/sys/class/net", link.Name())
        dev, err := os.Readlink(fullpath)
        if err != nil {
            continue 
        }
        if strings.HasPrefix(dev, "../../devices/virtual/") {
            continue 
        }
        
        systemConfig = &NICConfig{}
        err = systemConfig.GetConfig(link.Name())
        if err != nil {
            switch err.Error() {
            case "Parse nic link error.":
                continue 
            default:
                fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get nic info, %s\n", err))
            }   
        }
        nicInfo := *systemConfig.(*NICConfig)
        nicInfoFmt := nicInfo.GetInfoFmt()
        fmtInfo = append(fmtInfo, fmt.Sprintf("%d: ", index))
        fmtInfo = append(fmtInfo, nicInfoFmt)
        index  += 1
    }
 
    // filesystem
    fmtInfo = append(fmtInfo, "FileSystem Info:\n")
    lines, err := utils.ReadLines("/proc/mounts")
    if err != nil {
        fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get filesystem info, %s\n", err))
    }
    index = 1
    for _, line := range lines {
        systemConfig = &FileSystemConfig{}
        err = systemConfig.GetConfig(line)
        if err != nil {
            switch err.Error() {
            case "Fields < 4.":
                continue
            case "Doesn't judge this type.":
                continue
            default:
                fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get filesystem info, %s\n", err))
            }
        }
        fileSystemInfo    := *systemConfig.(*FileSystemConfig)
        fileSystemInfoFmt := fileSystemInfo.GetInfoFmt()
        fmtInfo = append(fmtInfo, fmt.Sprintf("%d: ", index))
        fmtInfo = append(fmtInfo, fileSystemInfoFmt)
        index  += 1
    } 

    // stroage
    fmtInfo = append(fmtInfo, "Stroage Info:\n")
    diskDevs, err := utils.DiskDev()
    if err != nil {
        fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get stroage info, %s\n", err))
    }
    index = 1
    for _, diskDev := range diskDevs {
        fullpath := path.Join("/sys/block", diskDev)
        dev, err := os.Readlink(fullpath)
        if err != nil {
            continue                       
        }
        if strings.HasPrefix(dev, "../../devices/virtual/") {
            continue 
        }

        deviceType, err := ioutil.ReadFile(path.Join(fullpath, "device", "type"))
        if err != nil {
            continue 
        } else if strings.HasPrefix(dev, "../devices/platform/floppy") || string(deviceType) == "5" {
            continue
        }

        systemConfig = &StroageConfig{}
        err = systemConfig.GetConfig(diskDev)
        if err != nil {
            switch err.Error() {
            case "Parse stroage link error.":
                continue 
            default:
                fmtInfo = append(fmtInfo, fmt.Sprintf("Can't get stroage info, %s\n", err))
            }
        }
        stroageInfo := *systemConfig.(*StroageConfig)
        stroageInfoFmt := stroageInfo.GetInfoFmt()
        fmtInfo = append(fmtInfo, fmt.Sprintf("%d: ", index))
        fmtInfo = append(fmtInfo, stroageInfoFmt)
        index  += 1 
    }
    return fmtInfo
}
