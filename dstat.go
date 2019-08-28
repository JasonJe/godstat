package main

import (
	"fmt"
	"encoding/json"

	stat "./stat"
)

func main() {
	sysStat := new(stat.SysStat)
	sysStat.CpuUtilization(200)
	sysStatJson, err := json.MarshalIndent(sysStat, "", "\t")
	if err != nil {
		fmt.Println(nil)
	}

	fmt.Println(string(sysStatJson))
}