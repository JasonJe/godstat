package main

import (
	stat "./stat"
)

func main() {
    sysStat := new(stat.SysStat)
    sysStat.Run(1000)
}
