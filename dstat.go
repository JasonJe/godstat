package main

import (
	stat "godstat/stat"
)

func main() {
    sysStat := new(stat.SysStat)
    sysStat.Run(1000)
}
