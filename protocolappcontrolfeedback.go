package sha

import (
	"encoding/binary"
)

type FeedbackAppControlPacket struct {
	Uid             uint64
	BedVersion      uint8
	ProtocolVersion uint8

	Result uint8
}

func (p *FeedbackAppControlPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xAC)
	gatewayid := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid, p.Uid)
	buf = append(buf, gatewayid[2:]...)
	buf = append(buf, 0xED)

	return buf
}

func NewLoginPakcet(Uid uint64, BedVersion uint8, ProtocolVersion uint8) {
	return &LoginPacket{
		Uid:             Uid,
		BedVersion:      BedVersion,
		ProtocolVersion: ProtocolVersion,
	}
}

func ParseLogin(buffer []byte, c *Conn) *LoginPacket {
	reader := bytes.NewReader(buffer)
	reader.Seek(3, 0)
	uid := make([]byte, 6)
	reader.Read(uid)
	gid := []byte{0, 0}
	gid = append(gid, uid...)
	bedid = binary.BigEndian.Uint64(gid)

	bedversion, _ := reader.ReadByte()
	protocolversion, _ := reader.ReadByte()

	NewGatewayHub().add(bedid, bedversion, protocolversion)
	c.uid = bedid
	c.SetStatus(ConnSuccess)
	NewConns().Add(c)

	return NewLoginPakcet(bedid, bedversion, protocolversion)
}
