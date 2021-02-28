package ethernet

import (
	"bytes"
	"encoding/binary"

	"github.com/gasugesu/lectcp/pkg/net"
)


type header struct {
	Dst Address
	Src Address
	Type net.EthernetType
}

type frame struct {
	header
	payload []byte
}

func parse(data []byte) (*frame, error) {
	f := frame{}
	buf := bytes.NewBuffer(data)
	if err:=binary.Read(buf, binary.BigEndian, &f.header); err != nil {
		return nil, err
	}
	f.payload = buf.Bytes()
	return &f, nil
}