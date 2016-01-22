package bed

import (
	"bytes"
	"encoding/binary"
)

type HandleControlFeedbackPacket struct {
	SerialNum uint32
}

func (p *HandleControlFeedbackPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xAA)
	buf = append(buf, 6)
	buf = append(buf, 0x02)
	buf = append(buf, 0x00)

	serialnum_byte := make([]byte, 4)
	binary.BigEndian.PutUint32(serialnum_byte, p.SerialNum)
	buf = append(buf, serialnum_byte...)
	sum := CheckSum(buf[2:], 6)
	buf = append(buf, sum)
	buf = append(buf, 0xED)

	return buf
}

func ParseHandleControlFeedback(buffer []byte) *HandleControlFeedbackPacket {
	reader := bytes.NewReader(buffer)
	reader.Seek(14, 0)

	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	return &HandleControlFeedbackPacket{
		SerialNum: serialnumber,
	}
}
