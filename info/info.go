package main

import (
    "fmt"
    "unsafe"
    "syscall"
    "strings"

    utils "../utils"
)

type SystemConfig interface {
    GetConfig() error
}

type KernelConfig struct {
    Release       string `json:"release"`
    Version       string `json:"version"`
    Architercture string `json:"architercture"`
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
    // unsafe.Pointer() 包含任意类型的地址
    // (*[65]byte)(unsafe.Pointer(&uname.Machine)) 将该地址装维 byte 数组
    // strings.TrimRight() 删除字符串右边指定字符
    kernelConfig.Architercture = strings.TrimRight(string((*[65]byte)(unsafe.Pointer(&uname.Machine))[:]), "\000")
    return nil
}

type OsConfig struct {
    Name          string `json:"name"`
    Vendor        string `json:"vendor"`
    Version       string `json:"version"`
    Release       string `json:"release"`
    Architercture string `json:"architercture"`
}

func (osConfig *OsConfig) GetConfig() error {
    osConfig.Name = "CentOS Linux 7 (Core)"
    osConfig.Vendor = "centos"
    osConfig.Version = "7"
    osConfig.Release = "7.2.1511"
    osConfig.Architercture = "amd64"
    return nil
}

func main() {
    var config SystemConfig 
    
    config = &KernelConfig{}
    config.GetConfig()
    fmt.Println(config)

    config = &OsConfig{}
    config.GetConfig()
    fmt.Println(config)
}
