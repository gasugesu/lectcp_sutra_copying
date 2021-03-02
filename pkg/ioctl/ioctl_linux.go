package ioctl

import (
	"bytes"
	"syscall"
	"unsafe"
)

func SIOCGIFINDEX(name string) (int32, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		// TODO: why not syscall.IFNAMSIZ
		name  [16]byte
		index int32
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	// TODO: Why is ifreq casted into uintptr(unsafe.Pointer())
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFINDEX, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return 0, errno
	}
	return ifreq.index, nil
}

// TODO: get flag
func SIOCGIFFLAGS(name string) (uint16, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return 0, errno
	}
	return ifreq.flags, nil
}

// TODO: set flag?
func SIOCSIFFLAGS(name string, flags uint16) error {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	ifreq.flags = flags
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return errno
	}
	return nil
}

type sockaddr struct {
	family uint16
	addr   [14]byte
}

// TODO: get hardware address?
func SIOCGIFHWAADDR(name string) ([]byte, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name [syscall.IFNAMSIZ]byte
		// TODO: where can i find 'sockaddr' definition
		addr sockaddr
		// TODO: what is pad? Why is this pad 8?
		_pad [8]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFHWADDR, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return nil, errno
	}
	return ifreq.addr.addr[:], nil
}

func TUNSETIFF(fd uintptr, name string, flags uint16) (string, error) {
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], []byte(name))
	ifreq.flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return "", errno
	}
	return string(ifreq.name[:bytes.IndexByte(ifreq.name[:], 0)]), nil
}
