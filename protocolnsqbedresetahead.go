package bed

import (
	"github.com/giskook/bed/pb"
	"github.com/golang/protobuf/proto"
)

type NsqBedResetAheadPacket struct {
	Uid          uint64
	SerialNumber uint32
}

func (p *NsqBedResetAheadPacket) Serialize() []byte {
	command := &Report.Command{
		Type: Report.Command_CMT_REPBEDRESET,
	}
	report := &Report.ControlReport{
		Tid:     p.Uid,
		Command: command,
	}

	data, _ := proto.Marshal(report)
	return data
}

func ParseNsqBedResetAhead(serialnum uint32) *NsqBedResetAheadPacket {
	return &NsqBedResetAheadPacket{
		SerialNumber: serialnum,
	}
}
