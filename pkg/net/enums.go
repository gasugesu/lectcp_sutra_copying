package net

type HardwareType uint16

const (
	HardwareLoopback     = 0x0000
	HardwareTypeEthernet = 0x0001
)

type EthernetType uint16

const (
	EthernetTypeIP   EthernetType = 0x8000
	EthernetTypeARP  EthernetType = 0x0806
	EthernetTypeIPv6 EthernetType = 0x86dd
)

func (t EthernetType) String() string {
	switch t {
	case EthernetTypeIP:
		return "IP"
	case EthernetTypeARP:
		return "ARP"
	case EthernetTypeIPv6:
		return "IPv6"
	default:
		return "Unknown"
	}
}

type ProtocolNumber uint8
