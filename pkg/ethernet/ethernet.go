package ethernet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/gasugesu/lectcp/pkg/net"
	"github.com/gasugesu/lectcp/pkg/raw"
)

const (
	headerSize     = 14
	trailerSize    = 0 // without FCS
	maxPayloadSize = 1500
	minPayloadSize = 46
	minFrameSize   = headerSize + minPayloadSize + trailerSize
	maxFrameSize   = headerSize + maxPayloadSize + trailerSize
)

type Device struct {
	raw  raw.Device
	addr Address
	mtu  int
}

var _ net.LinkDevice = &Device{} // interface check

func NewDevice(raw raw.Device) (*Device, error) {
	if raw == nil {
		return nil, fmt.Errorf("raw device is required")
	}
	addr := Address{}
	copy(addr[:], raw.Address())
	return &Device{
		raw:  raw,
		addr: addr,
		mtu:  maxPayloadSize,
	}, nil
}

func (d *Device) Address() net.HardwareAddress {
	return d.addr
}

func (d *Device) BroadCastAddress() net.HardwareAddress {
	return BroadCastAddress
}

func (d *Device) Close() {
	d.raw.Close()
}

func (d *Device) HeaderSize() int {
	return headerSize
}

func (d *Device) MTU() int {
	return d.mtu
}

func (d *Device) Name() string {
	return d.raw.Name()
}

func (d *Device) NeedARP() bool {
	return true
}

func (d *Device) Read(buf []byte) (int, error) {
	return d.raw.Read(buf)
}

func (d *Device) RxHandler(data []byte, callback net.LinkDeviceCallbackHandler) {
	f, err := parse(data)
	if err != nil {
		log.Panicln(err)
		return
	}
	if f.Dst != d.addr {
		if !f.Dst.isGroupAddress() {
			// other hst address
			return
		}
		if f.Dst != BroadCastAddress {
			// multicast frame unsupported
			return
		}
	}
	if f.Src == d.addr {
		// loopback frame

	}
	callback(d, f.Type, f.payload, f.Src, f.Dst)
}

func (d *Device) SetAddress(addr Address) {
	d.addr = addr
}

func (d *Device) Tx(Type net.EthernetType, data, dst []byte) error {
	hdr := header{
		Dst:  NewAddress(dst),
		Src:  d.addr,
		Type: Type,
	}
	f := bytes.NewBuffer(make([]byte, 0))
	binary.Write(f, binary.BigEndian, hdr)
	binary.Write(f, binary.BigEndian, data)
	if pad := minFrameSize - f.Len(); pad > 0 {
		binary.Write(f, binary.BigEndian, bytes.Repeat([]byte{byte(0)}, pad))
	}
	_, err := d.raw.Write(f.Bytes())
	return err
}

func (d *Device) Type() net.HardwareType {
	return net.HardwareTypeEthernet
}
