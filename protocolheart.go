package bed

import (
	"bytes"
	"encoding/binary"
)

type HeartPacket struct {
	Uid uint64
}

func (this *HeartPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xAA)
	buf = append(buf, 7)
	buf = append(buf, 0xFF)
	mac := make([]byte, 8)
	binary.BigEndian.PutUint64(mac, this.Uid)
	buf = append(buf, mac[2:]...)
	sum := CheckSum(buf[2:], 7)
	buf = append(buf, sum)
	buf = append(buf, 0xED)

	return buf
}

func ParseHeart(buffer []byte) *HeartPacket {
	reader := bytes.NewReader(buffer)
	reader.Seek(3, 0)
	uid := make([]byte, 6)
	reader.Read(uid)
	gid := []byte{0, 0}
	gid = append(gid, uid...)
	bedid := binary.BigEndian.Uint64(gid)

	return &HeartPacket{
		Uid: bedid,
	}
}
