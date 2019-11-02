package core

import (
    "strconv"
    "strings"

    utils "godstat/utils"
)

type AIOStat struct {
    Requests  int64 `json:"numRequests"`
}

func (aioStat *AIOStat) AIOTicker() error {
    lines, err := utils.ReadLines("/proc/sys/fs/aio-nr")
    if err != nil {
        return err
    }
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 1 {
            continue
        }
        tempRead, _ := strconv.ParseInt(fields[0], 10, 64)
        aioStat.Requests = tempRead
    }
    return nil
}

