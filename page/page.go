package page

import (
	"strings"
	"strconv"

	utils "../utils"
)

type PageStat struct {
	PageIn   int64  `json:"pageIn"`
	PageOut  int64 `json:"pageOut"`
}

func (pageStat *PageStat) PageTicker() error {
	filename := "/proc/vmstat"
	lines, err := utils.ReadLines(filename)

	if err != nil {
		return err
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		if strings.HasPrefix(fields[0], "pswpin") == true {
			pageStat.PageIn, _ = strconv.ParseInt(fields[1], 10, 64)
		}
		if strings.HasPrefix(fields[0], "pswpout") == true {
			pageStat.PageOut, _ = strconv.ParseInt(fields[1], 10, 64)
		}
	}
	
	return nil
}