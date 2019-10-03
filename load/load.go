package load 

import (
    "strings"
    "strconv"

    utils "godstat/utils"
)

type LoadStat struct {
    Load1    float64 `json:"load1"`
    Load5    float64 `json:"load5"`
    Load15   float64 `json:"load15"`
}

func (loadStat *LoadStat) LoadTicker() error {
    filename   := "/proc/loadavg"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return err
    }
    
    for _, line := range lines {
        fields  := strings.Fields(line)
        if len(fields) < 3 {
            continue 
        }
        
        load1,  _ := strconv.ParseFloat(fields[0], 64)
        load5,  _ := strconv.ParseFloat(fields[1], 64)
        load15, _ := strconv.ParseFloat(fields[2], 64)

        loadStat.Load1  = load1 
        loadStat.Load5  = load5 
        loadStat.Load15 = load15 
    }
    return nil 
}

