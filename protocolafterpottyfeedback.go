package bed

import (
	"bytes"
	"encoding/binary"
	"github.com/giskook/bed/pb"
	"github.com/golang/protobuf/proto"
)

type FeedbackAfterPottyPacket struct {
	Uid              uint64
	PottyType        uint8
	PottyTime        uint16
	PottyWeight      uint16
	WaterTemperature uint8
	CloudTemperature uint8
	SerialNumber     uint32
}

func (p *FeedbackAfterPottyPacket) Serialize() []byte {
	afterpotty := &Report.ToiletComplete{
		Style:  uint32(p.PottyType),
		Time:   uint32(p.PottyTime),
		Weight: uint32(p.PottyWeight),
		Water:  uint32(p.WaterTemperature),
		Wind:   uint32(p.CloudTemperature),
	}

	command := &Report.Command{
		Type:   Report.Command_CMT_REPTOILETCOMPLETE,
		Toilet: afterpotty,
	}

	report := &Report.ControlReport{
		Tid:          p.Uid,
		SerialNumber: p.SerialNumber,
		Command:      command,
	}

	data, _ := proto.Marshal(report)

	return data
}

func ParseAfterPottyFeedback(buffer []byte, c *Conn) *FeedbackAfterPottyPacket {
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

	return &FeedbackAfterPottyPacket{
		Uid:              c.GetBedID(),
		PottyType:        pottytype,
		PottyTime:        pottytime,
		PottyWeight:      pottyweight,
		WaterTemperature: watertemperature,
		CloudTemperature: cloudtemperature,
		SerialNumber:     serialnumber,
	}
}
