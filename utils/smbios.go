package utils 

import (
    "os"
    "errors"
    "strings"
    "strconv"
    "io/ioutil"
    "bytes"
    "syscall"
    "encoding/binary"
)

func checkSum(s []byte) (bool) {
    var sum byte
    for _, e := range s {
        sum += e
    }
    if sum == 0 {
        return true 
    } else {
        return false
    }
}

func structureTableAddrEFI(f *os.File) (int64, int, error) {
    lines, err := ReadLines("/sys/firmware/efi/systab")
    if err     != nil {
        return 0, 0, err
    }
    
    for _, line := range lines {
        fields  := strings.Split(line, "=")
        if len(fields) != 2 || fields[0] != "SMBIOS" {
            continue
        } else {
            addr, err := strconv.ParseInt(fields[1], 0, 64)
            if err    != nil {
                return 0, 0, err 
            }

            eps, err  := syscall.Mmap(int(f.Fd()), addr, 0x1f, syscall.PROT_READ, syscall.MAP_SHARED)
            if err    != nil {
                return 0, 0, err
            }
            defer syscall.Munmap(eps)

            if checkSum(eps) && checkSum(eps[0x10:]) && bytes.Equal(eps[0x10: 0x15], []byte("_DMI_")) {
                return int64(binary.LittleEndian.Uint32(eps[0x18: 0x18 + 4])), int(binary.LittleEndian.Uint16(eps[0x16: 0x16 + 2])), nil
            }
        }
    }
    return 0, 0, errors.New("SMBIOS entry point not found")
}

func structureTableAddr(f *os.File) (int64, int, error) {
    mem, err := syscall.Mmap(int(f.Fd()), 0xf0000, 0x10000, syscall.PROT_READ, syscall.MAP_SHARED)
    if err   != nil {
        return 0, 0, err
    }
    defer syscall.Munmap(mem)
    
    for i := range mem {
        if i > len(mem) - 0x1f {
            break
        }

        if i % 16 != 0 || !bytes.Equal(mem[i: i + 4], []byte("_SM")) {
            continue
        }
        
        eps := mem[i: i + 0x1f]
        if checkSum(eps) && checkSum(eps[0x10:]) && bytes.Equal(eps[0x10: 0x15], []byte("_DMI_")) {
            return int64(binary.LittleEndian.Uint32(eps[0x18: 0x18 + 4])), int(binary.LittleEndian.Uint16(eps[0x16: 0x16 + 2])), nil
        }
    }
    return 0, 0, errors.New("SMBIOS entry point not found")
}

func StructureTable() ([]byte, error) {
    f, err := os.Open("/dev/mem")
    if err != nil {
        mem, err := ioutil.ReadFile("/sys/firmware/dmi/tables/DMI")
        if err   != nil {
            return nil, err
        }
        return mem, err 
    }
    defer f.Close()

    addr, length, err := structureTableAddrEFI(f)
    if err != nil {
        if addr, length, err = structureTableAddr(f); err != nil {
            return nil, err
        }
    }

    align    := addr & (int64(os.Getpagesize()) - 1)
    mem, err := syscall.Mmap(int(f.Fd()), addr - align, length + int(align), syscall.PROT_READ, syscall.MAP_SHARED)
    if err   != nil {
        return nil, err
    }
    return mem[align:], nil
}
