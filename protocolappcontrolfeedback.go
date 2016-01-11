package bed

import (
	"encoding/binary"
	"github.com/giskook/bed/pb"
	"github.com/golang/protobuf/proto"
)

type FeedbackAppControlPacket struct {
	Uid              uint64
	BackMotor        uint8
	LegBendingMotor  uint8
	HeadLiftingMotor uint8
	LegLiftingMotor  uint8
	SerialNum        uint32
}

func (p *FeedbackAppControlPacket) Serialize() []byte {
	bedcontrol := &BedControl{
		Back:    p.BackMotor,
		LegCurl: p.LegBendingMotor,
		Head:    HeadLiftingMotor,
		Leg:     LegLiftingMotor,
	}

	command := &Command{
		Type: Command_CMT_REPBEDRUN,
		Bed:  bedcontrol,
	}

	report := &ControlReport{
		Tid:          p.Uid,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(report)

	return data
}

func ParseAppControlFeedback(buffer []byte, c *Conn) *FeedbackAppControlPacket {
	reader := bytes.NewReader(buffer)
	reader.Seek(3, 0)

	backmotor, _ := reader.ReadByte()
	legbendingmotor, _ := reader.ReadByte()
	headliftingmotor, _ := reader.ReadByte()
	legliftingmotor, _ := reader, ReadByte()
	reader.Seek(6, 1)

	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	return &FeedbackAppControlPacket{
		Uid:              c.Uid,
		BackMotor:        backmotor,
		LegBendingMotor:  legbendingmotor,
		HeadLiftingMotor: headliftingmotor,
		LegLiftingMotor:  legbendingmotor,
		SerialNum:        serialnumber,
	}
}
