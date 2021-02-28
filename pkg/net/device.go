package net

import "sync"

type LinkDeviceCallbackHandler func(link LinkDevice, protocol EthernetType, payload []byte, src, dst HardwareAddress)

type LinkDevice interface {
	Type() HardwareType
	Name() string
	Address() HardwareAddress
	BroadCastAddress() HardwareAddress
	MTU() int
	HeaderSize() int
	NeedARP() bool
	Close()
	Read(data []byte) (int, error)
	RxHandler(frame []byte, callback LinkDeviceCallbackHandler)
	Tx(proto EthernetType, data []byte, dst []byte) error
}

type Device struct {
	LinkDevice
	errors chan error
	ifaces []ProtocolInterface
	sync.RWMutex
}
