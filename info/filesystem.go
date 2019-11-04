package info 

import (
    "fmt"
    "syscall"
    "errors"
    "strings"

    utils "godstat/utils"
)

type FileSystemConfig struct {
    Filesystem string `json:"fileSystem"`
    Type       string `json:"type"`
    Size       uint64 `json:"size"`
    Used       uint64 `json:"used"`
    Available  uint64 `json:"available"`
    MountedOn  string `json:"mountedOn"`
}

func (fileSystemConfig *FileSystemConfig) GetConfig(args ...interface{}) error {
    if line, ok := args[0].(string); ok {
        fields := strings.Fields(line)
        if len(fields) < 4 {
            return errors.New("Fields < 4.") 
        }
        
        if fields[1] == "" || fields[1] == "swap"|| fields[1] == "/dev/shm" || fields[1] == "/dev/pts" || strings.HasPrefix(fields[1], "/proc") || strings.HasPrefix(fields[1], "/sys") || strings.HasPrefix(fields[0], "/dev/loop") || fields[0] == "cgroup" || fields[0] == "securityfs" || fields[0] == "mqueue" || fields[0] == "hugetlbfs" || fields[0] == "sunrpc" || fields[0] == "lxcfs" || fields[0] == "shm" || fields[0] == "tmpfs" || fields[0] == "overlay" || fields[0] == "nsfs" {
            return errors.New("Doesn't judge this type.") 
        }

        fileSystemConfig.Filesystem = fields[0]
        fileSystemConfig.MountedOn  = fields[1] 
        fileSystemConfig.Type       = fields[2]
        
        fs  := syscall.Statfs_t{}
        err := syscall.Statfs(fields[1], &fs)
        if err != nil {
            return err 
        }
        fileSystemConfig.Size = fs.Blocks * uint64(fs.Bsize)
        fileSystemConfig.Available = fs.Bfree * uint64(fs.Bsize)
        fileSystemConfig.Used = fileSystemConfig.Size - fileSystemConfig.Available
        return nil
   }
   return errors.New("Unkown")
}

func (fileSystemConfig *FileSystemConfig) GetInfoFmt() string {
    fileSystemInfoFmt := fmt.Sprintf("\t%-40s\n\t%-40s\n\t%-40s\n\t%-40s\n", 
                                    "FileSystem: " + fileSystemConfig.Filesystem, 
                                    "Type: " + fileSystemConfig.Type + "\tSize: " + utils.ByteCountSI(int64(fileSystemConfig.Size)),
                                    "Used: " + utils.ByteCountSI(int64(fileSystemConfig.Used)) + "\tAvailable: " + utils.ByteCountSI(int64(fileSystemConfig.Available)),
                                    "MountedOn: " + fileSystemConfig.MountedOn)
    return fileSystemInfoFmt 
}
