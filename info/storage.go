package info 

import (
    "os"
    "path"
    "strconv"
    "errors"
    "strings"
    "io/ioutil"
    
    utils "godstat/utils"
)

type StroageConfig struct {
    Name   string `json:"stroageName"`
    Driver string `json:"stroageDriver"`
    Vendor string `json:"stroageVendor"`
    Model  string `json:"stroageModel"`
    Serial string `json:"stroageSerial"`
    Size   uint64 `json:"stroageSize"`
}

func (stroageConfig *StroageConfig) GetConfig(args ...interface{}) error {
    if name, ok := args[0].(string); ok {
        stroageConfig.Name = name
        fullpath := path.Join("/sys/block", name)
        modelRead,  err := ioutil.ReadFile(path.Join(fullpath, "device", "model"))
        if err != nil {
            return err 
        }
        stroageConfig.Model = strings.TrimSpace(string(modelRead))
        
        serialRead, err := ioutil.ReadFile(path.Join(fullpath, "dev"))
        if err != nil {
            return err
        } else if string(serialRead) != "" {
            lines, err := utils.ReadLines(path.Join("/run/udev/data", "b" + strings.TrimSpace(string(serialRead))))
            if err != nil {
                return err 
            }
            for _, line := range lines {
                fields := strings.Split(line, "=")
                if len(fields) == 2 {
                    if fields[0] == "E:ID_SERIAL_SHORT" {
                        stroageConfig.Serial = fields[1]
                        break 
                    }
                }
            } 
        }
        
        if driverRead, err := os.Readlink(path.Join(fullpath, "device", "driver")); err == nil {
            stroageConfig.Driver = path.Base(driverRead)
        }

        vendorRead, err := ioutil.ReadFile(path.Join(fullpath, "device", "vendor"))
        if err != nil {
            return err
        } else if !strings.HasPrefix(string(vendorRead), "0x"){
            stroageConfig.Vendor = strings.TrimSpace(string(vendorRead))
        }

        sizeRead, err := ioutil.ReadFile(path.Join(fullpath, "size"))
        if err != nil {
            return err
        } else if !strings.HasPrefix(string(vendorRead), "0x"){
            size, err := strconv.ParseInt(strings.TrimSpace(string(sizeRead)), 10, 64)
            if err != nil {
                return err
            }
            stroageConfig.Size = uint64(size) / 1953125 
        }
        
    }
    return errors.New("Unknow")
}
