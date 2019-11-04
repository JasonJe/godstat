package info 

import (
    "os"
    "fmt"
    "unsafe"
    "syscall"
    "strings"
    "io/ioutil"

    utils "godstat/utils"
)

type OsConfig struct {
    Name          string `json:"osName"`
    Vendor        string `json:"osVendor"`
    Version       string `json:"osVersion"`
    Release       string `json:"osRelease"`
    Architercture string `json:"osArchitercture"`
    HostName      string `json:"hostName"`
    TimeZone      string `json:"timeZone"`
}

func (osConfig *OsConfig) GetConfig(args ...interface{}) error {
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
            osConfig.Name    = strings.Replace(fields[1], "\"", "", -1)
        case "ID":
            osConfig.Vendor  = fields[1]
        case "VERSION_ID":
            osConfig.Version = strings.Replace(fields[1], "\"", "", -1)
        case "PRETTY_NAME":
            osConfig.Release = strings.Replace(fields[1], "\"", "", -1)
        }
    }
    
    hostRead, err := ioutil.ReadFile("/proc/sys/kernel/hostname")
    if err != nil {
        return err               
    }
    osConfig.HostName = strings.TrimSpace(string(hostRead))

    timeZoneRead, err := ioutil.ReadFile("/etc/timezone")
    if err != nil {
        return err               
    }
    osConfig.TimeZone = strings.TrimSpace(string(timeZoneRead))
    return nil
}

func (osConfig *OsConfig) GetInfoFmt() string {
    osInfoFmt := fmt.Sprintf("OS Info:\n\tName: %s\n\tVendor: %s\n\tVersion: %s\n\tRelease: %s\n\tArchitercture: %s\n\tHostName: %s\n\tTimeZone: %s\n", 
                            osConfig.Name, 
                            osConfig.Vendor, 
                            osConfig.Version, 
                            osConfig.Release, 
                            osConfig.Architercture, 
                            osConfig.HostName, 
                            osConfig.TimeZone)
    return osInfoFmt 
}
