package main

import (
	"time"
	"runtime"

	stat "./stat"
)



func main() {
	sysStat := new(stat.SysStat)

	ticker := time.NewTicker(time.Millisecond * 200)
	cpusStat, _ := sysStat.CpuTimes()
	<- ticker.C
	cpusStat2, _ := sysStat.CpuTimes()

	for i := 0; i < runtime.NumCPU() + 1; i++ {

		cpuName := cpusStat[i].CPU

		user2 := cpusStat2[i].User
		user := cpusStat[i].User

		nice2 := cpusStat2[i].Nice
		nice := cpusStat[i].Nice

		system2 := cpusStat2[i].System
		system := cpusStat[i].System

		idle2 := cpusStat2[i].Idle
		idle := cpusStat[i].Idle

		iowait2 := cpusStat2[i].Iowait
		iowait := cpusStat[i].Iowait

		irq2 := cpusStat2[i].Irq
		irq := cpusStat[i].Irq

		softirq2 := cpusStat2[i].Softirq
		softirq := cpusStat[i].Softirq

		steal2 := cpusStat2[i].Steal
		steal := cpusStat[i].Steal

		guest2 := cpusStat2[i].Guest
		guest := cpusStat[i].Guest

		guestNice2 := cpusStat2[i].GuestNice
		guestNice := cpusStat[i].GuestNice

		stolen2 := cpusStat2[i].Stolen
		stolen := cpusStat[i].Stolen

		cpu := user + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice + stolen
		cpu2 := user2 + nice2 + system2 + idle2 + iowait2 + irq2 + softirq2 + steal2 + guest2 + guestNice2 + stolen2

		puser := 100 * (user2 - user) / (cpu2 - cpu)
		pnice := 100 * (nice2 - nice) / (cpu2 - cpu)
		psystem := 100 * (system2 - system) / (cpu2 - cpu)
		pidle := 100 * (idle2 - idle) / (cpu2 - cpu)
		piowait := 100 * (iowait2 - iowait) / (cpu2 - cpu)
		pirq := 100 * (irq2 - irq) / (cpu2 - cpu)
		psoftirq := 100 * (softirq2 - softirq) / (cpu2 - cpu)
		psteal := 100 * (steal2 - steal) / (cpu2 - cpu)
		pguest := 100 * (guest2 - guest) / (cpu2 - cpu)

		
	}
}