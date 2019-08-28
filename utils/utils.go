package utils

import (
	"os"
	"bufio"
	"strings"
)

func ReadLines(filename string) ([]string, error) {
	return readLinesOffsetN(filename, 0, -1)
}

func readLinesOffsetN(filename string, offset, limit int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < (limit + offset) || limit < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < limit {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil
}