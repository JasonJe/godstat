package core

import (
    "strings"
    "strconv"

    utils "godstat/utils"
)

type SystemStat struct {
    Interrupt     float64 `json:"interrupt"`
    ContextSwitch float64 `json:"contextSwitch"`
}

func (systemStat *SystemStat) SystemTicker() error {
    filename := "/proc/stat"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return err
    }
    
    var interrupt, contextSwitch float64
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 2 {
            continue 
        }

        if fields[0] == "intr" {
            interrupt, _ = strconv.ParseFloat(fields[1], 64)
        }
        if fields[0] == "ctxt" {
            contextSwitch, _ = strconv.ParseFloat(fields[1], 64)
        }
    }
    systemStat.Interrupt     = interrupt 
    systemStat.ContextSwitch = contextSwitch 

    return nil
}
