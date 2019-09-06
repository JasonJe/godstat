package utils

import (
	"time"
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"reflect"
	"path/filepath"
)

type FormatTime time.Time

func (this FormatTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

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
            if err == io.EOF {
                break 
            }
			return nil, err
		}
		if i < limit {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	
    r.Reset(r) // *** 关键，这里的 f 是 *os.File，Reset 将 r 的底层 Reader 重新指定为 r，同时丢弃缓存中的所有数据，复位
    // f.close()
	return ret, nil
}

func GetDiskDev() ([]string, error) {
    files, _ := filepath.Glob("/sys/block/*")
    var baseNames []string
    for _, file := range files {
         baseName := strings.Split(file, "/")
         baseNames = append(baseNames, baseName[len(baseName) - 1])
     }
    
    filename := "/proc/diskstats"
    lines, err := readLinesOffsetN(filename, 0, -1)
    if err != nil {
        return nil, err
    }
    for _, line := range lines { 
        fields := strings.Fields(line)
        
        if len(fields) < 13 {
            continue
        }
        unUsingDstat := []string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"}
        if (reflect.DeepEqual(fields[3:], unUsingDstat)) {
            continue
        }
        if stringsContains(baseNames, fields[2]) {
            continue 
        }
        baseNames = append(baseNames, fields[2])
    }
    return baseNames, nil
}

func stringsContains(array []string, val string) (bool) {
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            return true  
        }
    }
    return false 
}
