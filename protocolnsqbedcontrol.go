package bed

import (
	"encoding/binary"
	"github.com/giskook/bed/pb"
)

type NsqBedControlPacket struct {
	SerialNumber     uint32
	BackMotor        uint8
	LegBendingMotor  uint8
	HeadLiftingMotor uint8
	LegLiftingMotor  uint8
}

func (p *NsqBedControlPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xAA)
	buf = append(buf, 15)
	buf = append(buf, 0x01)
	buf = append(buf, p.BackMotor)
	buf = append(buf, p.LegBendingMotor)
	buf = append(buf, p.HeadLiftingMotor)
	buf = append(buf, p.LegBendingMotor)
	buf = append(buf, 0)
	buf = append(buf, 0)
	buf = append(buf, 0)
	buf = append(buf, 0)
	buf = append(buf, 0)
	buf = append(buf, 0)
	serialnum_byte := make([]byte, 4)
	binary.BigEndian.PutUint32(serialnum_byte, p.SerialNumber)
	buf = append(buf, serialnum_byte...)
	sum := CheckSum(buf[2:], 15)
	buf = append(buf, sum)
	buf = append(buf, 0xED)

	return buf
}

func ParseNsqBedControl(serialnum uint32, command *Report.Command) *NsqBedControlPacket {
	return &NsqBedControlPacket{
		SerialNumber:     serialnum,
		BackMotor:        uint8(command.Bed.Back),
		LegBendingMotor:  uint8(command.Bed.LegCurl),
		HeadLiftingMotor: uint8(command.Bed.Head),
		LegLiftingMotor:  uint8(command.Bed.Leg),
	}
}
