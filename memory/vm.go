package memory  

import (
    "strings"
    "strconv"

    utils "godstat/utils"
)

type VMStat struct {
    PgMajFault int64 `json:"PageMajorFault"`
    PgFault      int64 `json:"PageFault"`
    PgAlloc      int64 `json:"PageAlloc"`
    PgFree       int64 `json:"PageFree"`
}

func (vmStat *VMStat) VMTicker() error {
    lines, err := utils.ReadLines("/proc/vmstat")
    if err != nil {
        return err
    }
    for _, line := range lines {
        fields  := strings.Fields(line)
        if len(fields) < 2 {
            continue 
        }
        valueRead, err := strconv.ParseInt(fields[1], 10, 64)
        if err != nil {
            return err
        }

        switch {
        case fields[0] == "pgmajfault":
            vmStat.PgMajFault = valueRead 
        case fields[0] == "pgfault":
            vmStat.PgFault    = valueRead 
        case fields[0] == "pgfree":
            vmStat.PgFree     = valueRead 
        case strings.HasPrefix(fields[0], "pgalloc_"):
            vmStat.PgAlloc    = vmStat.PgAlloc + valueRead 
        }
    }
    return nil
}

