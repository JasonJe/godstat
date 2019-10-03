package info 

import (
    "syscall"
    "errors"
    "strings"
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
            return errors.New("Unkown") 
        }
        
        if fields[1] == "" || fields[1] == "swap"|| fields[1] == "/dev/shm" || fields[1] == "/dev/pts" || strings.HasPrefix(fields[1], "/proc") || strings.HasPrefix(fields[1], "/sys") || strings.HasPrefix(fields[0], "/dev/loop") || fields[0] == "cgroup" || fields[0] == "securityfs" || fields[0] == "mqueue" || fields[0] == "hugetlbfs" || fields[0] == "sunrpc" || fields[0] == "lxcfs" || fields[0] == "shm" || fields[0] == "tmpfs" || fields[0] == "overlay" || fields[0] == "nsfs" {
            return errors.New("Unkown") 
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

