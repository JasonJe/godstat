package filesystem

import (
	"fmt"
	"strings"
	"strconv"

	utils "godstat/utils"
)

type FileSystemStat struct {
	UsingFileHandle int64 `json:"usingFileHandle"`
	UsingInode     int64 `json:"usingInode"`
}

func (fileSystemStat *FileSystemStat) FileSystemTicker() error {
	lines, err := utils.ReadLines("/proc/sys/fs/file-nr")
	if err != nil {
		return err
	}
	for _, line := range lines {
		fields  := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		filesHandle, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return err
		}
		fileSystemStat.UsingFileHandle = filesHandle
	}

	lines, err = utils.ReadLines("/proc/sys/fs/inode-nr")
	if err != nil {
		return err
	}
	for _, line := range lines {
		fields  := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		usedInode, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return err
		}
		freeInode, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return err
		}
		fileSystemStat.UsingInode = usedInode - freeInode
	}
	return nil
}