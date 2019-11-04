package core

import (
	"strings"
	"strconv"

	utils "godstat/utils"
)

type MemoryStat struct {
	Total        float64  `json:"total"`
	Free         float64  `json:"free"`
	Buffers      float64  `json:"buffers"`
	Cached       float64  `json:"cached"`
	Active       float64  `json:"active"`
	Inactive     float64  `json:"inactive"`
	Available    float64  `json:"available"`
	Used         float64  `json:"used"`
	UsedPercent  float64  `json:"usedPercent"`
    SReclaimable float64  `json:"sReclaimable"`
    Shmem        float64  `json:"shmem"`
}

func (memoryStat *MemoryStat) MemoryTicker() error {
	filename := "/proc/meminfo"
	lines, err := utils.ReadLines(filename)
	if err != nil {
		return err
	}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
	
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		t, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		switch key {
		case "MemTotal":
			memoryStat.Total        = t * 1024.0
		case "MemFree":
			memoryStat.Free         = t * 1024.0
		case "Buffers":
			memoryStat.Buffers      = t * 1024.0
		case "Cached":
			memoryStat.Cached       = t * 1024.0
		case "Active":
			memoryStat.Active       = t * 1024.0
		case "Inactive":
			memoryStat.Inactive     = t * 1024.0
		case "SReclaimable":
		    memoryStat.SReclaimable = t * 1024.0
		case "Shmem":
		    memoryStat.Shmem        = t * 1024.0
		}

	}
	memoryStat.Available   = memoryStat.Free + memoryStat.Buffers + memoryStat.Cached
	memoryStat.Used        = memoryStat.Total - memoryStat.Free - memoryStat.Buffers - memoryStat.Cached - memoryStat.SReclaimable - memoryStat.Shmem
	memoryStat.UsedPercent = memoryStat.Used / memoryStat.Total * 100.0
	return nil
}
