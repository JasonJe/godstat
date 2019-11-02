package core

import (
    "strconv"
    "strings"
    
    utils "godstat/utils"
)

type SocketStat struct {
    Total int64 `json:"socketTotal"`
    TCP   int64 `json:"numTCP"`
    UDP   int64 `json:"numUDP"`
    RAW   int64 `json:"numRAW"`
    FRAG  int64 `json:"numFRAG"`
    Other int64 `json:"numOther"`
}

func (socketStat *SocketStat) SocketTicker() error {
    filename   := "/proc/net/sockstat"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return err 
    }

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 3 {
            continue 
        }
        num, err := strconv.ParseInt(fields[2], 10, 64)
        if err != nil {
            return err 
        }
        switch fields[0] {
        case "sockets:": 
            socketStat.Total = num
        case "TCP:":
            socketStat.TCP   = num
        case "UDP:":
            socketStat.UDP   = num
        case "RAW:":
            socketStat.RAW   = num
        case "FRAG:":
            socketStat.FRAG  = num
        }
    }
    socketStat.Other = socketStat.Total - socketStat.TCP - socketStat.UDP - socketStat.RAW - socketStat.FRAG 
    return nil 
}

