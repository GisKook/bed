package bed

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

type NsqBedPottyPacket struct {
	SerialNumber uint32
}

func (p *NsqBedPottyPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xAA)
	buf = append(buf, 7)
	buf = append(buf, 0x03)
	buf = append(buf, 0)
	buf = append(buf, 0)
	serialnum_byte := make([]byte, 4)
	binary.BigEndian.PutUint32(serialnum_byte, p.SerialNumber)
	buf = append(buf, serialnum_byte)
	sum := CheckSum(buf[2:], 7)
	buf = append(buf, sum)
	buf = append(buf, 0xED)

	return buf
}

func ParseNsqPotty(serialnum uint32) *NsqBedPottyPacket {
	return &NsqBedPottyPacket{
		SerialNumber: serialnum,
	}
}
