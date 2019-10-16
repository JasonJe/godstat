package lock

import (
    "strings"

    utils "godstat/utils"
)

type LockStat struct {
    Posix int64 `json:"Posix"`
    Flock int64 `json:"Flock"`
    Read  int64 `json:"Read"`
    Write int64 `json:"Write"`
}

func (lockStat *LockStat) LockTicker() error {
    lines, err := utils.ReadLines("/proc/locks")
    if err != nil {
        return err
    }
    lockStat.Posix = 0
    lockStat.Flock = 0
    lockStat.Read  = 0
    lockStat.Write = 0
    for _, line := range lines {
        fields  := strings.Fields(line)
        if len(fields) < 4 {
            continue
        }
        switch {
        case fields[1] == "POSIX":
            lockStat.Posix = lockStat.Posix + 1
        case fields[1] == "FLOCK":
            lockStat.Flock = lockStat.Flock + 1
        }
        switch {
        case fields[3] == "READ":
            lockStat.Read  = lockStat.Read  + 1
        case fields[3] == "WRITE":
            lockStat.Write = lockStat.Write + 1
        }
    }
    return nil
}

