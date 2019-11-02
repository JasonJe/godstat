package core

import (
    "strconv"
    "strings"

    utils "godstat/utils"
)

type ZoneStat struct {
    DMA32Free  int64 `json:"DMA32Free"`
    DMA32High  int64 `json:"DMA32High"`
    NormalFree int64 `json:"NormalFree"`
    NormalHigh int64 `json:"NormalHigh"`
}

func (zoneStat *ZoneStat) ZoneTicker() error {
    lines, err  := utils.ReadLines("/proc/zoneinfo")
    if err != nil {
        return err
    }
    for index, line := range lines {
        fields  := strings.Fields(line)
        if len(fields) < 2 {
            continue
        }
        if strings.HasPrefix(fields[0], "Node") {
            detail := lines[index: index + 9]
            switch {
            case fields[3] == "DMA32":
                freeFields    := strings.Fields(detail[1])
                freeRead, err := strconv.ParseInt(freeFields[2], 10, 64)
                if err != nil {
                    return err 
                }
                highFields    := strings.Fields(detail[4])
                highRead, err := strconv.ParseInt(highFields[1], 10, 64)
                if err != nil {
                    return err
                }
                zoneStat.DMA32Free = freeRead 
                zoneStat.DMA32High = highRead
            case fields[3] == "Normal":
                freeFields    := strings.Fields(detail[1])
                freeRead, err := strconv.ParseInt(freeFields[2], 10, 64)
                if err != nil {
                    return err 
                }
                highFields    := strings.Fields(detail[4])
                highRead, err := strconv.ParseInt(highFields[1], 10, 64)
                if err != nil {
                    return err
                }
                zoneStat.NormalFree = freeRead 
                zoneStat.NormalHigh = highRead
            }
        }
    }
    return nil
}

