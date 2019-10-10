package socket

import (
    utils "godstat/utils"
)

type RawSocketStat struct {
    NumSockets int64 `json:"numSockets"`
}

func (rawSocketStat *RawSocketStat) RawSocketTicker() error {
    filename   := "/proc/net/raw"
    lines, err := utils.ReadLines(filename)
    if err != nil {
        return err 
    }
    rawSocketStat.NumSockets = int64(len(lines) - 1)
    return nil
}
