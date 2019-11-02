package core

import (
    "strconv"
    "strings"
    
    utils "godstat/utils"
)

type ProcStat struct {
    Running   float64 `json:"processRunning"`
    Blocked   float64 `json:"processBlocked"`
    Processes float64 `json:"processes"`
}

func (procStat *ProcStat) ProcTicker() error {
    filename   := "/proc/stat"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return err
    }

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 2 {
            continue 
        }
        if fields[0] == "processes" {
            tempRead, err := strconv.ParseFloat(fields[1], 64)
            if err != nil {
                return err
            }
            procStat.Processes = tempRead 
        } else if fields[0] == "procs_running" {
            tempRead, err := strconv.ParseFloat(fields[1], 64)
            if err != nil {
                return err
            }
            procStat.Running = tempRead - 1.0
        } else if fields[0] == "procs_blocked" {
            tempRead, err := strconv.ParseFloat(fields[1], 64)

            if err != nil {
                return err
            }
            procStat.Blocked = tempRead
        }
    }
    return nil
}


