package main

import (
	_"fmt"
	_"time"
	// "encoding/json"

	stat "./stat"
	// utils "./utils"
)

func main() {
	// sysStat := new(stat.SysStat)

	// sysStat.CpuUtilization(200)
	// sysStat.MemoryInfo()
	// sysStat.Paging(200)
	
	// sysStat.DateTime = utils.FormatTime(time.Now())

	// sysStatJson, err := json.MarshalIndent(sysStat, "", "\t")
	// if err != nil {
	// 	fmt.Println(nil)
	// }

	// fmt.Println(string(sysStatJson))

//	for {
//		sysStat := new(stat.SysStat)
//		// sysStat.Paging(200)
//		// fmt.Println(sysStat.PageIn, sysStat.PageOut)
//		sysStat.Disk(200)
//        // fmt.Println(sysStat.DiskList)
//		time.Sleep(time.Second)
//	}
    
    sysStat := new(stat.SysStat)
    sysStat.Run(1000)

}
