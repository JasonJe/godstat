package info 

import (
    "bytes"
    "syscall"
    "encoding/binary"    
    
    utils "godstat/utils"
)

var (
    memTypes = [30]string{
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
}

func (memoryConfig *MemoryConfig) GetConfig(args ...interface{}) error { 
    mem, err := utils.StructureTable()
    if err != nil {
        return err 
    } 
    defer syscall.Munmap(mem) // mmap 将一个文件或者其它对象映射进内存, munmap 解除内存映射

    var memSizeAlt uint 
    for p := 0; p < len(mem) - 1; {
        recType := mem[p]
        recLen  := mem[p + 1]

        switch recType {
        case 17: 
            size := uint(binary.LittleEndian.Uint16(mem[p + 0x0c: p + 0x0c + 2]))
            if size == 0 || size == 0xffff || size & 0x8000 == 0x8000 {
                break
            }
            if size == 0x7fff {
                if recLen >= 0x20 {
                    size = uint(binary.LittleEndian.Uint32(mem[p + 0x1c: p + 0x1c + 4]))
                } else {
                    break
                }
            }
            memoryConfig.Size += uint16(size)
            if index := int(mem[p + 0x12]); index >= 1 && index <= len(memTypes) {
                memoryConfig.Type = memTypes[index - 1]
            }

            if recLen >= 0x17 {
                if speed := uint(binary.LittleEndian.Uint16(mem[p + 0x15: p + 0x15 + 2])); speed != 0 {
                    memoryConfig.Speed = uint16(speed)
                }
            }
        
        case 19:
            start := uint(binary.LittleEndian.Uint32(mem[p + 0x04: p + 0x04 + 4]))
            end   := uint(binary.LittleEndian.Uint32(mem[p + 0x08: p + 0x08 + 4]))

            if start == 0xffffffff && end == 0xffffffff {
                if recLen >= 0x1f {
                    start64 := binary.LittleEndian.Uint64(mem[p + 0x0f: p + 0x0f + 8])
                    end64   := binary.LittleEndian.Uint64(mem[p + 0x17: p + 0x17 + 8])
                    memSizeAlt += uint((end64 - start64 + 1) / (1024 * 1024))
                }
            } else {
                memSizeAlt += (end - start + 1) / 1024
            }
        
        case 127 :
            break
        }

        for p+= int(recLen); p < len(mem) - 1; {
            if bytes.Equal(mem[p:p+2], []byte{0, 0}) {
                p += 2
                break
            }
            p++
        }
    }
    if memoryConfig.Size == 0 && memSizeAlt > 0 {
        memoryConfig.Type = "DRAM"
        memoryConfig.Size = uint16(memSizeAlt)
    }
    return nil    
}
