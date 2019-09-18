package disk

import (
    "regexp"
    "strings"
    "strconv"
    
    utils "../utils"
)

type DiskStat struct {
    Name        string  `json:"name"`
    Read        float64 `json:"read"`
    Write       float64 `json:"write"`
}

func (diskStat *DiskStat) getDiskStat(tempRead, tempWrite float64) {
    diskStat.Read  = diskStat.Read  + tempRead 
    diskStat.Write = diskStat.Write + tempWrite
}

func DiskTicker() ([]DiskStat, error) {
    filename := "/proc/diskstats"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return nil, err 
    }
    
    devNames, err := utils.DiskDev()
    if err != nil {
        return nil, err
    }

    diskList := []DiskStat{}
    for _, diskDev := range devNames {
        diskStat := DiskStat{diskDev, 0.0, 0.0}
        diskList = append(diskList, diskStat)
    }
    
    totalDiskStat := DiskStat{"total", 0.0, 0.0}
    for _, line := range lines {
        fields := strings.Fields(line)
        
        devIndex := utils.StringsContains(devNames, fields[2]) 
        if devIndex == -1 {
            continue 
        } else {
            tempRead,  _ := strconv.ParseFloat(fields[5], 64)
            tempWrite, _ := strconv.ParseFloat(fields[9], 64)
            
            diskStat := diskList[devIndex]
            diskStat.getDiskStat(tempRead, tempWrite)
            diskList[devIndex] = diskStat
            
            if isMatch, _ := regexp.MatchString(`^([hsv]d[a-z]+\d+|cciss/c\d+d\d+p\d+|dm-\d+|md\d+|mmcblk\d+p\d0|VxVM\d+)$`, fields[2]); !isMatch {
               totalDiskStat.Read  = totalDiskStat.Read  + tempRead 
               totalDiskStat.Write = totalDiskStat.Write + tempWrite
            }
        }
    }
    diskList = append(diskList, totalDiskStat)
    return diskList, nil
}
