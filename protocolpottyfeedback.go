package bed

import (
	"encoding/binary"
	"github.com/giskook/bed/pb"
	"github.com/golang/protobuf/proto"
)

type FeedbackPottyPacket struct {
	Uid       uint64
	SerialNum uint32
	CmdType   uint8
}

func (p *FeedbackPottyPacket) Serialize() []byte {
	bedcontrol := &BedControl{
		Back:    p.BackMotor,
		LegCurl: p.LegBendingMotor,
		Head:    HeadLiftingMotor,
		Leg:     LegLiftingMotor,
	}

	command := nil
	if p.CmdType == AppPottyFeedback {
		command = &Command{
			Type: Command_CMT_REPTOILET,
			Bed:  bedcontrol,
		}
	} else {
		command = &Command{
			Type: Command_CMT_REPMANUALTOILET,
			Bed:  bedcontrol,
		}
	}

	report := &ControlReport{
		Tid:          p.Uid,
		SerialNumber: p.SerialNum,
		Command:      command,
	}

	data, _ := proto.Marshal(report)

	return data
}

func ParsePottyFeedback(buffer []byte, c *Conn, cmdtype uint8) *FeedbackPotty {
	reader := bytes.NewReader(buffer)
	reader.Seek(5, 0)
	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	return &FeedbackPotty{
		Uid:       c.Uid,
		SerialNum: serialnumber,
	}
}
