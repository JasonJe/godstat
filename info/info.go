package main

import (
    "fmt"
    "os"
    "path" 
    "strings"
    "io/ioutil"
    
    "github.com/digitalocean/go-smbios/smbios"

    utils "../utils"
)

var SS [](*smbios.Structure)

type SystemConfig interface {
    GetConfig(args ...interface{}) error
}

func init() {
    // SMBIOS 
    rc, _, err := smbios.Stream()
    if err != nil {
        panic(err)
    }
    // Be sure to close the stream!
    defer rc.Close()

    // Decode SMBIOS structures from the stream.
    d := smbios.NewDecoder(rc)
    SS, err = d.Decode()
    if err != nil {
        panic(err)
    }
}

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
    
    for _, s := range SS {
        if byte(s.Header.Type) == 4 {
            config = &CpuConfig{}
            config.GetConfig(s)
            fmt.Println(config)    
        }
    }

    for _, s := range SS {
        if byte(s.Header.Type) == 17 {
            config = &MemoryConfig{}
            e := config.GetConfig(s)
            if e != nil {
                continue
            }
            fmt.Println(config)
        }
    }

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
        config = &NICConfig{}
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
        config  = &StroageConfig{}
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
        config = &FileSystemConfig{}
        err := config.GetConfig(line)
        if err != nil {
            continue
        }
        fmt.Println(config)
    }
}
