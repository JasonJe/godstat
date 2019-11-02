package core 

import (
    "strings"

    utils "godstat/utils"
)

type UDPStat struct {
    Established int64 `json:"numEstablished"`
    Listen      int64 `json:"numListen"`
} 

func (udpStat *UDPStat) selectType(value string) {
    switch value {
    case "01":
        (*udpStat).Established = (*udpStat).Established + 1
    case "07":
        (*udpStat).Listen      = (*udpStat).Listen      + 1
    }
}

func (udpStat *UDPStat) UDPTicker() error {
    lines, err := utils.ReadLines("/proc/net/udp")
    if err != nil {
        return err 
    }
    for _, line := range lines {
        fields  := strings.Fields(line)
        udpStat.selectType(fields[3])
    }

    lines2 , err := utils.ReadLines("/proc/net/udp6")
    if err != nil {
        return err 
    }
    for _, line := range lines2 {
        fields  := strings.Fields(line)
        udpStat.selectType(fields[3])
    }
    return nil
}
