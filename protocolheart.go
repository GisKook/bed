package sha

import (
	"encoding/binary"
)

type HeartPacket struct {
	Uid uint64
}

func (this *HeartPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xCE)
	buf = append(buf, 0x00)
	buf = append(buf, 0x0B)
	buf = append(buf, 0x80)
	buf = append(buf, 0x02)
	gatewayid := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid, this.Uid)
	buf = append(buf, gatewayid[2:]...)
	buf = append(buf, CheckSum(buf, uint16(len(buf))))
	buf = append(buf, 0xCE)

	return buf
}

func NewHeartPacket(Uid uint64) *HeartPacket {
	return &HeartPacket{
		Uid: Uid,
	}
}

func ParseHeart(buffer []byte) *HeartPacket {
	gatewayid, reader := GetGatewayID(buffer)
	reader.ReadByte()
	reader.ReadByte()

	return NewHeartPacket(gatewayid)
}
