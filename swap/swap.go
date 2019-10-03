package swap

import (
    "strings"
    "strconv"

    utils "godstat/utils"
)

type SwapStat struct {
    Name  string  `json:"name"`
    Used  float64 `json:"used"`
    Free  float64 `json:"free"`
}

func SwapTicker() ([]SwapStat, error) {
    filename   := "/proc/swaps"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return nil, err
    }
    
    swapList := []SwapStat{}
    totalSwapStat := SwapStat{"name", 0.0, 0.0}
    for _, line := range lines {
        fields := strings.Fields(line)
        if (len(fields) < 5) || (fields[0] == "Filename") {
            continue 
        }
        
        used, _  := strconv.ParseFloat(fields[3], 64)
        total, _ := strconv.ParseFloat(fields[2], 64)
        
        swapStat := SwapStat{fields[0], used * 1024.0, (total - used) * 1024.0}
        
        totalSwapStat.Used = totalSwapStat.Used + used * 1024.0
        totalSwapStat.Free = totalSwapStat.Free + (total - used) * 1024.0
        
        swapList = append(swapList, swapStat)
    }
    swapList = append(swapList, totalSwapStat)
    return swapList, nil
}
