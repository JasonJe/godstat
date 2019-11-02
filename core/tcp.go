package core

import (
    "strings"

    utils "godstat/utils"
)

type TCPStat struct {
    Established  int64 `json:"numEstablished"`
    SynSent      int64 `json:"numSynSent"`
    SynRecv      int64 `json:"numSynRecv"`
    FinWait1     int64 `json:"numFinWait1"`
    FinWait2     int64 `json:"numFinWait2"`
    TimeWait     int64 `json:"numTimeWait"`
    Close        int64 `json:"numClose"`
    CloseWait    int64 `json:"numCloseWait"`
    LastAck      int64 `json:"numLastAck"`
    Listen       int64 `json:"numListen"`
    Closing      int64 `json:"numClosing"`
}

func (tcpStat *TCPStat) selectType(value string) {
    /* 01: established 
       02: syn_sent 
       03: syn_recv 
       04: fin_wait1 
       05: fin_wait2 
       06: time_wait 
       07: close 
       08: close_wait 
       09: last_ack 
       0A: listen 
       0B: closing
    */

    switch value {
    case "01":
        (*tcpStat).Established = (*tcpStat).Established + 1
    case "02":
        (*tcpStat).SynSent     = (*tcpStat).SynSent     + 1
    case "03":
        (*tcpStat).SynRecv     = (*tcpStat).SynRecv     + 1
    case "04":
        (*tcpStat).FinWait1    = (*tcpStat).FinWait1    + 1
    case "05":
        (*tcpStat).FinWait2    = (*tcpStat).FinWait2    + 1
    case "06":
        (*tcpStat).TimeWait    = (*tcpStat).TimeWait    + 1
    case "07":
        (*tcpStat).Close       = (*tcpStat).Close       + 1
    case "08":
        (*tcpStat).CloseWait   = (*tcpStat).CloseWait   + 1
    case "09":
        (*tcpStat).LastAck     = (*tcpStat).LastAck     + 1
    case "0A":
        (*tcpStat).Listen      = (*tcpStat).Listen      + 1
    case "0B":
        (*tcpStat).Closing     = (*tcpStat).Closing     + 1
    }
}

func (tcpStat *TCPStat) TCPTicker() error { 
    lines, err  := utils.ReadLines("/proc/net/tcp")
    if err != nil {
        return err 
    }
    for _, line := range lines {
        fields  := strings.Fields(line)
        if len(fields) < 12 {
            continue 
        }
        tcpStat.selectType(fields[3])
    }
    
    lines2, err := utils.ReadLines("/proc/net/tcp6")
    if err != nil {
        return err 
    }
    for _, line := range lines2 {
        fields := strings.Fields(line)
        if len(fields) < 12 {
            continue 
        }
        tcpStat.selectType(fields[3])
    }
    return nil 
}
