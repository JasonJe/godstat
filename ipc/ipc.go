package ipc 

import (
    "path"

    utils "godstat/utils"
)

type IPCStat struct {
    MessageQueue int64 `json:"IPCMessageQueue"`
    Semaphore    int64 `json:"IPCSemaphore"`
    SharedMemory int64 `json:"IPCSharedMemory"`
}

func (ipcStat *IPCStat) IPCTicker() error {
    for _, name    := range [3]string{"msg", "sem", "shm"} {
        filename   := path.Join("/proc/sysvipc/", name)
        lines, err := utils.ReadLines(filename)
        if err != nil {
            return err
        }
        switch name {
        case "msg":
            ipcStat.MessageQueue = int64(len(lines) - 1)
        case "sem":
            ipcStat.Semaphore    = int64(len(lines) - 1)
        case "shm":
            ipcStat.SharedMemory = int64(len(lines) - 1)
        }
    }
    return nil
}
