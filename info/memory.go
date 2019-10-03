package info 

import (
    "bytes"
    "errors"
    "encoding/binary"    
    
    "github.com/digitalocean/go-smbios/smbios"
)

var (
    memoryTypes = [30]string{
		"Other",
		"Unknown",
		"DRAM",
		"EDRAM",
		"VRAM",
		"SRAM",
		"RAM",
		"ROM",
		"FLASH",
		"EEPROM",
		"FEPROM",
		"EPROM",
		"CDRAM",
		"3DRAM",
		"SDRAM",
		"SGRAM",
		"RDRAM",
		"DDR",
		"DDR2",
		"DDR2 FB-DIMM",
		"Reserved1",
		"Reserved2",
		"Reserved3",
		"DDR3",
		"FBD2",
		"DDR4",
		"LPDDR",
		"LPDDR2",
		"LPDDR3",
		"LPDDR4"}
)

type MemoryConfig struct {
    Size         uint16 `json:"memorySize"`
    Type         string `json:"memoryType"`
    Speed        uint16 `json:"memorySpeed"`
    SerialNumber string `json:"memorySerialNumber"`
}

func (memoryConfig *MemoryConfig) GetConfig(args ...interface{}) error { 
    // 安全类型断言，失败value则为类型T对应的零值
    // value,ok := expression.(T) 
    // 非安全类型断言，失败时会panic()
    // value := expression.(T)
    if s, ok := args[0].(*smbios.Structure); ok { // 类型断言 
        if 0x14 > len(s.Formatted) - 1 {
            return errors.New("Unknown")               
        }
        index := s.Formatted[0x14]
        if index == 0 {
            return errors.New("Unknown")
        }
        memoryConfig.SerialNumber = s.Strings[index - 1]
        var sizeU uint16 
        binary.Read(bytes.NewBuffer(s.Formatted[0x08: 0x0A][0:2]), binary.LittleEndian, &sizeU)
        memoryConfig.Size = sizeU 
        memoryConfig.Type = memoryTypes[s.Formatted[0xE]]
        var speedU uint16 
        binary.Read(bytes.NewBuffer(s.Formatted[0x11: 0x13][0:2]), binary.LittleEndian, &speedU)
        memoryConfig.Speed = speedU 
        return nil 
    }
    return errors.New("Unknown")
}
