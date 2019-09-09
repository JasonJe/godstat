package disk

import (
    "regexp"
    "strings"
    "strconv"
    "reflect"

    utils "../utils"
)

type DiskStat struct {
    Name        string  `json:"name"`
    Read        float64 `json:"read"`
    Write       float64 `json:"write"`
}

func DiskTicker(totalDiskStat *DiskStat) (map[string]DiskStat, error) {
    filename := "/proc/diskstats"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return nil, err 
    }
    diskList := make(map[string]DiskStat)

    diskList["total"] = *totalDiskStat

    for _, line := range lines { 
        fields := strings.Fields(line)
        
        if len(fields) < 13 {
                continue
        }
        if fields[5] == "0" && fields[9] == "0" {
            continue
        }
        unUsingDstat := []string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"}
        if (reflect.DeepEqual(fields[3:], unUsingDstat)) {
            continue
        }
        
        tempRead,  _ := strconv.ParseFloat(fields[5], 64)
        tempWrite, _ := strconv.ParseFloat(fields[9], 64)
        
        if isMatch, _ := regexp.MatchString(`^([hsv]d[a-z]+\d+|cciss/c\d+d\d+p\d+|dm-\d+|md\d+|mmcblk\d+p\d0|VxVM\d+)$`, fields[2]); !isMatch {
           (*totalDiskStat).Read  = (*totalDiskStat).Read  + tempRead 
           (*totalDiskStat).Write = (*totalDiskStat).Write + tempWrite
        }
        
        if _, ok := diskList[fields[2]]; ok {
            diskStat := DiskStat{fields[2], diskList[fields[2]].Read  + tempRead, diskList[fields[2]].Write + tempWrite}
            diskList[fields[2]] = diskStat
        } else {
            diskStat := DiskStat{fields[2], 0 + tempRead, 0 + tempWrite}
            diskList[fields[2]] = diskStat 
        }
    }

    return diskList, nil
} 
