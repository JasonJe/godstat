package core 

import (
    "regexp"
    "strings"
    "strconv"

    utils "godstat/utils"
)

type IOStat struct {
    Name        string  `json:"name"`
    Read        float64 `json:"read"`
    Write       float64 `json:"write"`
}

func (ioStat *IOStat) getIOStat(tempRead, tempWrite float64) {
    ioStat.Read  = ioStat.Read  + tempRead 
    ioStat.Write = ioStat.Write + tempWrite 
}

func IOTicker() ([]IOStat, error) {
    filename := "/proc/diskstats"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return nil, err 
    }
    
    devNames, err := utils.DiskDev()
    if err != nil {
        return nil, err 
    } 
    
    ioList := []IOStat{}
    for _, diskDev := range devNames {
        ioStat := IOStat{diskDev, 0.0, 0.0}
        ioList  = append(ioList, ioStat)
    }
    
    totalIOStat := IOStat{"total", 0.0, 0.0}
    for _, line := range lines {
        fields := strings.Fields(line)
        
        devIndex := utils.StringsContains(devNames, fields[2]) 
        if devIndex == -1 {
            continue 
        } else {
            tempRead,  _ := strconv.ParseFloat(fields[3], 64)
            tempWrite, _ := strconv.ParseFloat(fields[7], 64)
            
            ioStat := ioList[devIndex]
            ioStat.getIOStat(tempRead, tempWrite)
            ioList[devIndex] = ioStat
            
            if isMatch, _ := regexp.MatchString(`^([hsv]d[a-z]+\d+|cciss/c\d+d\d+p\d+|dm-\d+|md\d+|mmcblk\d+p\d0|VxVM\d+)$`, fields[2]); !isMatch {
               totalIOStat.Read  = totalIOStat.Read  + tempRead 
               totalIOStat.Write = totalIOStat.Write + tempWrite
            } 
        }
    }
    ioList  = append(ioList, totalIOStat)
    return ioList, nil 
}

