package core

import (
    "strings"

    utils "godstat/utils"
)

type UnixSocketStat struct {
    DataGram    int64 `json:"numDatagram"`
    Stream      int64 `json:"numStream"`
    Established int64 `json:"numEstablished"`
    Listen      int64 `json:"numListen"`
}

func (unixSocketStat *UnixSocketStat) UnixSocketTicker() error {
    lines, err := utils.ReadLines("/proc/net/unix")
    if err != nil {
        return err 
    }

    for _, line := range lines {
        fields  := strings.Fields(line)
        switch fields[4] {
        case "0002":
            (*unixSocketStat).DataGram = (*unixSocketStat).DataGram + 1
        case "0001":
            (*unixSocketStat).Stream   = (*unixSocketStat).Stream   + 1
            switch fields[5] {
            case "01":
                (*unixSocketStat).Listen      = (*unixSocketStat).Listen + 1
            case "03":
                (*unixSocketStat).Established = (*unixSocketStat).Established + 1
            }
        }
    }
    return nil
}

