package main

import (
    "fmt"
    "bytes"
    "syscall"
    "unsafe"
    "strings"
)

const (
    IFNAMSIZ          = 16
    PERMADDR_LEN      = 32
    SIOCETHTOOL       = 0x8946 // Linux SIOCETHTOOL ioctl 操作，可以用来进行网络设备的统计信息、驱动程序相关信息检索
    SIOCGIFADDR       = 0x8915 // 
    ETHTOOL_GDRVINFO  = 0x00000003 // 获取驱动信息指令
    ETHTOOL_GPERMADDR = 0x00000020 // 获取永久地址
    ETHTOOL_GSET      = 0x00000001 // 获取配置信息
)

type ethtool struct { 
    fd       int 
    cmd      uint32
    driver   string 
    macAddr  string
    ip       string
    maxSpeed uint 
    port     string
}

// 
type ethtoolDriverInfo struct {
    cmd          uint32 
    driver       [32]byte 
    version      [32]byte 
    fw_version   [32]byte 
    bus_info     [32]byte 
    erom_version [32]byte 
    reserved2    [12]byte 
    n_priv_flags uint32 
    n_stats      uint32 
    testinfo_len uint32  
    eedump_len   uint32 
    regdump_len  uint32
}

type ethtoolPermAddr struct {
    cmd  uint32 
    size uint32 
    data [PERMADDR_LEN]byte
}

type ethtoolCmd struct {
    cmd              uint32 
    supported        uint32
    advertising      uint32 
    speed            uint16 
    duplex           uint8 
    phy_address      uint8 
    transceiver      uint8 
    autoneg          uint8 
    mdio_support     uint8 
    maxtxpkt         uint32 
    maxrxpkt         uint32 
    speed_hi         uint16 
    eth_tp_mdix      uint8 
    eth_tp_mdix_ctrl uint8 
    lp_advertising   uint32
    reserved         [2]uint32
}

type ifre struct {
    name [IFNAMSIZ]byte 
    data uintptr 
}

type NICConfig struct {
    Name       string     `json:"nicName"`
    Driver     string     `json:"nicDriver"`
    MACAddress string     `json:"macAddress"`
    Ports      string     `json:"nicPorts"`
    Speed      uint       `json:"nicSpeed"`
    IP         string     `json:"nicIP"`
}

func newEthtool() (*ethtool, error) {
    // syscall.AF_INET，表示服务器之间的网络通信
    // syscall.SOCK_DGRAM, 基于UDP的socket通信，应用层socket。
    // syscall.IPPROTO_IP 接收任何的IP数据包
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
    if err != nil {
        return nil, err
    }
    
    return &ethtool{
        fd: int(fd)}, nil 
}

func (eth *ethtool) close() {
    syscall.Close(eth.fd)
}

func (eth *ethtool) ioctl(intf string, data uintptr) error {
    // syscall.SYS_IOCTL 建立一个 socket，得到一个 fd，然后在此 fd 上执行 ioctl 的各种操作
    var name [IFNAMSIZ]byte 
    copy(name[:], []byte(intf))

    ifr := ifre {
        name: name,
        data: data}
    
    // int ioctl(int d, int request, ...)
    // 第一个参数是打开的设备文件的文件描述符，通常是open系统调用的返回值；第二个参数request是可以自定义的请求号；第三个参数可以是一个指针，指向一段用户态内存，用来传递参数，也可以是一个整形数据。
    _, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(eth.fd), SIOCETHTOOL, uintptr(unsafe.Pointer(&ifr)))
    if ep != 0 {
        return syscall.Errno(ep)
    }
    return nil
}

func (eth *ethtool) driverName(intf string) (error) {
    driver := ethtoolDriverInfo{
        cmd: ETHTOOL_GDRVINFO}
    if err := eth.ioctl(intf, uintptr(unsafe.Pointer(&driver))); err != nil {
        return err 
    }
    eth.driver = string(bytes.Trim(driver.driver[:], "\x00"))
    return nil 
}

func (eth *ethtool) permAddr(intf string) (error) {
    addr := ethtoolPermAddr{
        cmd : ETHTOOL_GPERMADDR,
        size: PERMADDR_LEN}
    if err := eth.ioctl(intf, uintptr(unsafe.Pointer(&addr))); err != nil {
        return err 
    }
    
    eth.macAddr = fmt.Sprintf("%x:%x:%x:%x:%x:%x",
        addr.data[0:1],
        addr.data[1:2],
        addr.data[2:3],
        addr.data[3:4],
        addr.data[4:5],
        addr.data[5:6])
    return nil 
}

func (eth *ethtool) settingInfo(intf string) (error) {
    cmd := ethtoolCmd{
        cmd: ETHTOOL_GSET}
    if err := eth.ioctl(intf, uintptr(unsafe.Pointer( & cmd))); err != nil {
        return err 
    }
    
    switch {
    case cmd.supported & 0x78000000 > 0:
        eth.maxSpeed = 56000
    case cmd.supported & 0x07800000 > 0:
        eth.maxSpeed = 40000
    case cmd.supported & 0x00600000 > 0:
        eth.maxSpeed = 20000
    case cmd.supported & 0x001c1000 > 0:
        eth.maxSpeed = 10000
    case cmd.supported & 0x00008000 > 0:
        eth.maxSpeed = 2500
    case cmd.supported & 0x00020030 > 0:
        eth.maxSpeed = 1000
    case cmd.supported & 0x0000000c > 0:
        eth.maxSpeed = 100
    case cmd.supported & 0x00000003 > 0:
        eth.maxSpeed = 10
    }
    
    for i, p := range []string{"tp", "aui", "mii", "fibre", "bnc"} {
        if cmd.supported & (1 << uint(i + 7)) > 0 {
            eth.port += p + "/"
        }
    }

    eth.port = strings.TrimRight(eth.port, "/")
    return nil 
}

func (eth *ethtool) inetAddr(intf string) (error) {
    var ifreqbuf [40]byte 
    for i := 0; i < 40; i++ {
        ifreqbuf[i] = 0            
    }
    for i := 0; i < len(intf); i++ {
        ifreqbuf[i] = intf[i]
    }
    _, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(eth.fd), SIOCGIFADDR, uintptr(unsafe.Pointer(&ifreqbuf)))
    if ep != 0 {
        return syscall.Errno(ep)
    }
    eth.ip = fmt.Sprintf("%d.%d.%d.%d", ifreqbuf[20], ifreqbuf[21], ifreqbuf[22], ifreqbuf[23])
    return nil 
}

func (nicConfig *NICConfig) GetConfig(args ...interface{}) error {
    eth, err := newEthtool()
    if err != nil {
        return err   
    }

    err = eth.driverName(nicConfig.Name)
    if err != nil {
        return err
    }

    err = eth.permAddr(nicConfig.Name)
    if err != nil {
        return err
    }

    err = eth.inetAddr(nicConfig.Name)
    if err != nil {
        return err
    }

    err = eth.settingInfo(nicConfig.Name)
    if err != nil {
        return err
    }
    
    nicConfig.Driver = eth.driver 
    nicConfig.MACAddress = eth.macAddr 
    nicConfig.Ports = eth.port 
    nicConfig.Speed = eth.maxSpeed
    nicConfig.IP = eth.ip
   
    return nil
}
