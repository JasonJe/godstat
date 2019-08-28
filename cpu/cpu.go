package cpu

import (
	"strings"
	"strconv"
	"errors"
)

type CpuStat struct {
	CPU       string
	User      float64
	System    float64
	Idle      float64
	Nice      float64
	Iowait    float64
	Irq       float64
	Softirq   float64
	Steal     float64
	Guest     float64
	GuestNice float64
	Stolen    float64
}

func (cpuStat *CpuStat) GetCpuTimes(line string) error {
	fields := strings.Fields(line)

	if strings.HasPrefix(fields[0], "cpu") == false {
		return errors.New("not contain cpu")
	}

	cpu := fields[0]
	user, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return err
	}
	nice, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return err
	}
	system, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return err
	}
	idle, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return err
	}
	iowait, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return err
	}
	irq, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return err
	}
	softirq, err := strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return err
	}
	stolen, err := strconv.ParseFloat(fields[8], 64)
	if err != nil {
		return err
	}
	
	cpuTick := float64(100)
	cpuStat.CPU = cpu
	cpuStat.User = float64(user) / cpuTick
	cpuStat.System = float64(system) / cpuTick
	cpuStat.Idle = float64(idle) / cpuTick
	cpuStat.Nice = float64(nice) / cpuTick
	cpuStat.Iowait = float64(iowait) / cpuTick
	cpuStat.Irq = float64(irq) / cpuTick
	cpuStat.Softirq = float64(softirq) / cpuTick
	cpuStat.Stolen = float64(stolen) / cpuTick
	
	if len(fields) > 9 { // Linux >= 2.6.11
		steal, err := strconv.ParseFloat(fields[9], 64)
		if err != nil {
			return err
		}
		cpuStat.Steal = float64(steal)
	}
	if len(fields) > 10 { // Linux >= 2.6.24
		guest, err := strconv.ParseFloat(fields[10], 64)
		if err != nil {
			return err
		}
		cpuStat.Guest = float64(guest)
	}
	if len(fields) > 11 { // Linux >= 3.2.0
		guestNice, err := strconv.ParseFloat(fields[11], 64)
		if err != nil {
			return err
		}
		cpuStat.GuestNice = float64(guestNice)
	}
	return nil
}

