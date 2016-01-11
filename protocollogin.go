package bed

import (
	"encoding/binary"
)

var (
	Infrared      uint8 = 0
	DoorMagnetic  uint8 = 1
	WarningButton uint8 = 2
)

type LoginPacket struct {
	Uid             uint64
	BedVersion      uint8
	ProtocolVersion uint8

	Result uint8
}

func (p *LoginPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xAA)
	buf = append(buf, 7)
	buf = append(buf, 0)
	mac := make([]byte, 8)
	binary.BigEndian.PutUint64(mac, p.Uid)
	buf = append(buf, mac[2:]...)
	sum := CheckSum(buf[2:], 7)
	buf = append(buf, sum)
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
