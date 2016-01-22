package bed

import (
	"bytes"
	"encoding/binary"
)

type AfterPottyPacket struct {
	PottyType        uint8
	PottyTime        uint16
	PottyWeight      uint16
	WaterTemperature uint8
	CloudTemperature uint8
	SerialNumber     uint32
}

func (p *AfterPottyPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0xAA)
	buf = append(buf, 12)
	buf = append(buf, 0x05)
	buf = append(buf, p.PottyType)
	buf = append(buf, p.PottyType)

	pottytime_byte := make([]byte, 2)
	binary.BigEndian.PutUint16(pottytime_byte, p.PottyTime)
	buf = append(buf, pottytime_byte...)

	pottyweight_byte := make([]byte, 2)
	binary.BigEndian.PutUint16(pottyweight_byte, p.PottyWeight)
	buf = append(buf, pottyweight_byte...)
	buf = append(buf, p.WaterTemperature)
	buf = append(buf, p.CloudTemperature)

	serialnum_byte := make([]byte, 4)
	binary.BigEndian.PutUint32(serialnum_byte, p.SerialNumber)
	buf = append(buf, serialnum_byte...)
	sum := CheckSum(buf[2:], 12)
	buf = append(buf, sum)
	buf = append(buf, 0xED)

	return buf
}

func ParseAfterPottyTobedFeedback(buffer []byte) *AfterPottyPacket {
	reader := bytes.NewReader(buffer)
	reader.Seek(3, 0)

	pottytype, _ := reader.ReadByte()
	pottytime_byte := make([]byte, 2)
	reader.Read(pottytime_byte)
	pottytime := binary.BigEndian.Uint16(pottytime_byte)

	pottyweight_byte := make([]byte, 2)
	reader.Read(pottyweight_byte)
	pottyweight := binary.BigEndian.Uint16(pottyweight_byte)

	watertemperature, _ := reader.ReadByte()
	cloudtemperature, _ := reader.ReadByte()

	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	return &AfterPottyPacket{
		PottyType:        pottytype,
		PottyTime:        pottytime,
		PottyWeight:      pottyweight,
		WaterTemperature: watertemperature,
		CloudTemperature: cloudtemperature,
		SerialNumber:     serialnumber,
	}
}
