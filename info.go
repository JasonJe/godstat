package main

import (
   "fmt"
   "os"
   "path" 
   "strings"
   "io/ioutil"
    
    info  "godstat/info"
    utils "godstat/utils"
)

func main() {
    var config info.SystemConfig 
    
    config = &info.KernelConfig{}
    config.GetConfig()
    fmt.Println(config)

    config = &info.OsConfig{}
    err := config.GetConfig()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(config)
    
     config = &info.CpuConfig{}
     config.GetConfig()
     fmt.Println(config)    

    config = &info.MemoryConfig{}
    e := config.GetConfig()
    if e != nil {
        panic(e)
    }
    fmt.Println(config)

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
        config = &info.NICConfig{}
        config.GetConfig(link.Name())
        fmt.Println(config)
    }
    
    devNames, err := utils.DiskDev()
    if err != nil {
        panic(err)
    }

    for _, diskDev := range devNames {
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
        config  = &info.StroageConfig{}
        err = config.GetConfig(diskDev)
        if err != nil {
            continue
        }
        fmt.Println(config)
    }

    lines, err := utils.ReadLines("/proc/mounts")
    if err != nil {
        panic(err)
    }
    for _, line := range lines {
        config = &info.FileSystemConfig{}
        err := config.GetConfig(line)
        if err != nil {
            continue
        }
        fmt.Println(config)
    }
}
