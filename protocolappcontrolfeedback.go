package bed

import (
	"bytes"
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
	CmdType          uint8
}

func (p *FeedbackAppControlPacket) Serialize() []byte {
	bedcontrol := &Report.BedControl{
		Back:    uint32(p.BackMotor),
		LegCurl: uint32(p.LegBendingMotor),
		Head:    uint32(p.HeadLiftingMotor),
		Leg:     uint32(p.LegLiftingMotor),
	}

	var command *Report.Command
	if p.CmdType == AppControlFeedback {
		command = &Report.Command{
			Type: Report.Command_CMT_REPBEDRUN,
			Bed:  bedcontrol,
		}
	} else {
		command = &Report.Command{
			Type: Report.Command_CMT_REPMANUALBEDRUN,
			Bed:  bedcontrol,
		}
	}

	report := &Report.ControlReport{
		Tid:          p.Uid,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(report)

	return data
}

func ParseAppControlFeedback(buffer []byte, c *Conn, cmdtype uint8) *FeedbackAppControlPacket {
	reader := bytes.NewReader(buffer)
	reader.Seek(3, 0)

	backmotor, _ := reader.ReadByte()
	legbendingmotor, _ := reader.ReadByte()
	headliftingmotor, _ := reader.ReadByte()
	legliftingmotor, _ := reader.ReadByte()
	reader.Seek(6, 1)

	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	return &FeedbackAppControlPacket{
		Uid:              c.GetBedID(),
		BackMotor:        backmotor,
		LegBendingMotor:  legbendingmotor,
		HeadLiftingMotor: headliftingmotor,
		LegLiftingMotor:  legliftingmotor,
		SerialNum:        serialnumber,
		CmdType:          cmdtype,
	}
}
