package core

import (
    "regexp"
    "strings"
    "strconv"

    utils "godstat/utils"
)

type NetStat struct {
    Name        string  `json:"name"`
    Recv        float64 `json:"read"`
    Send        float64 `json:"write"`
}

func (netStat *NetStat) getNetStat(tempRecv, tempSend float64) {
    netStat.Recv = netStat.Recv + tempRecv 
    netStat.Send = netStat.Send + tempSend
}

func NetTicker() ([]NetStat, error) {
    filename := "/proc/net/dev"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return nil, err
    }
    lines = lines[2:] // **
    
    netList := []NetStat{}
    totalNetStat := NetStat{"total", 0.0, 0.0}
    for _, line := range lines {
        fields := strings.Fields(line)
        
        if len(fields) < 17 {
            continue
        }
        if (fields[2] == "0" && fields[10] == "0") {
            continue
        }
        if strings.Contains(fields[0], "face") {
            continue 
        }

        devName     := strings.Replace(fields[0], ":", "", 1)
        tempRecv, _ := strconv.ParseFloat(fields[1], 64)
        tempSend, _ := strconv.ParseFloat(fields[9], 64)
        
        netStat := NetStat{devName, tempRecv, tempSend}
        netList = append(netList, netStat)
        
        if isMatch, _ := regexp.MatchString(`^(lo|bond\d+|face|.+\.\d+)$`, devName); !isMatch {
            totalNetStat.Recv = totalNetStat.Recv + tempRecv 
            totalNetStat.Send = totalNetStat.Send + tempSend
        }
    }
    netList = append(netList, totalNetStat)
    return netList, nil
}
