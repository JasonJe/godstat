package info 

import (
    "fmt"
    "syscall"

    utils "godstat/utils"
)

type KernelConfig struct {
    Release       string `json:"kernelRelease"`
    Version       string `json:"kernelVersion"`
}

func (kernelConfig *KernelConfig) GetConfig(args ...interface{}) error {
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

func (kernelConfig *KernelConfig) GetInfoFmt() string {
    kernelInfoFmt := fmt.Sprintf("Kernel Info:\n\tRelease: %s\n\tVersion: %s\n", kernelConfig.Release, kernelConfig.Version)
    return kernelInfoFmt 
}
