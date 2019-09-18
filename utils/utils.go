package utils

import (
	"time"
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"errors"
	"reflect"
	"path/filepath"
)

type FormatTime time.Time

func (this FormatTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func ByteCountSI(b int64) string {
    const unit = 1024
    if b < unit {
        return fmt.Sprintf("%dB", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit 
        exp++                            
    }
    return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "kMGTPE"[exp])
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
	return ret, nil
}

func StringsContains(array []string, val string) int {
    for i := 0; i < len(array); i++ {
        if array[i] == val {
            return i  
        }
    }
    return -1 
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
        if StringsContains(baseNames, fields[2]) != -1 {
            continue 
        }
        baseNames = append(baseNames, fields[2])
    }
    return baseNames, nil
}

func DiskBaseName(disk string) (string, error) {
    if ok := strings.HasPrefix(disk, "/dev/"); ok {
        if _, err := os.Stat(disk); err != nil {
            return "", err 
        } else {
            diskInfo, err := filepath.EvalSymlinks(disk)
            if err != nil {
                return "", err
            } else {
                return strings.Replace(diskInfo, "/dev/", "", 1), nil
            }
        } 
    } else {
        return "", errors.New("disk does not exist.")
    }
}

func DiskDev() ([]string, error) {
    var devNames []string

    files, _ := filepath.Glob("/sys/block/*") 
    for _, file := range files {
         devName := strings.Split(file, "/")
         devNames = append(devNames, devName[len(devName) - 1])
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
        if StringsContains(devNames, fields[2]) != -1 {
            continue 
        }
        devNames = append(devNames, fields[2])
    }
    return devNames, nil
}

func NetDev() ([]string, error) {
    filename := "/proc/net/dev"
    lines, err := readLinesOffsetN(filename, 0, -1)
    if err != nil {
        return nil, err
    }
    
    var devNames []string
    for _, line := range lines { 
        fields := strings.Fields(line)
        
        if len(fields) < 17 {
            continue
        }
        if (fields[2] == "0" && fields[10] == "0") {
            continue
        }
        if strings.Contains(fields[0], "face") {
            continue 
        }

        devNames = append(devNames, strings.Replace(fields[0], ":", "", 1))
    }
    return devNames, nil
}
