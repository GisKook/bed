package bed

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

type NsqBedResetAheadPacket struct {
	Uid          uint64
	SerialNumber uint32
}

func (p *NsqBedResetAheadPacket) Serialize() []byte {
	command := &Command{
		Type: Command_CMT_REPBEDRESET,
	}
	report := &ControlReport{
		Tid:     p.Uid,
		Command: command,
	}

	data, _ := proto.Marshal(report)
	return data
}

func ParseNsqBedResetAhead(serialnum uint32) *NsqBedResetAheadPacket {
	return &NsqBedResetPacket{
		SerialNumber: serialnum,
	}
}
