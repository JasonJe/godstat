package net

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
    devNames, err := utils.NetDev()  
    if err != nil {
        return nil, err
    }
    
    netList := []NetStat{} 
    for _, netDev := range devNames {
        netStat := NetStat{netDev, 0.0, 0.0}
        netList = append(netList, netStat)
    }
    
    totalNetStat := NetStat{"total", 0.0, 0.0}
    for index, netStat := range netList {
        fields := strings.Fields(lines[index])
        
        tempRecv, _ := strconv.ParseFloat(fields[1], 64)
        tempSend, _ := strconv.ParseFloat(fields[9], 64)

        if utils.StringsContains(devNames, strings.Replace(fields[0], ":", "", 1)) != -1 {
            netStat.getNetStat(tempRecv, tempSend)
            netList[index] = netStat // 除了指针、map、slice、chan等引用类型外，所有传参都是传值，都是一个副本/拷贝 
        }
        if isMatch, _ := regexp.MatchString(`^(lo|bond\d+|face|.+\.\d+)$`, strings.Replace(fields[0], ":", "", 1)); !isMatch {
            totalNetStat.Recv = totalNetStat.Recv + tempRecv 
            totalNetStat.Send = totalNetStat.Send + tempSend
        }
    }
    netList = append(netList, totalNetStat)
    return netList, nil
}
