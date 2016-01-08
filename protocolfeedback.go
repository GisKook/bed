package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
)

type FeedbackPacket struct {
	GatewayID    uint64
	SerialNumber uint32
	Result       uint8
}

func (p *FeedbackPacket) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REPOPFEEDBACK,
		Paras: para,
	}
	feedback := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNumber,
		Command:      command,
	}

	data, _ := proto.Marshal(feedback)

	return data
}

func ParseFeedback(buffer []byte) *FeedbackPacket {
	gatewayid, reader := GetGatewayID(buffer)
	serialnumber_byte := make([]byte, 4)
	reader.Read(serialnumber_byte)
	serialnumber := binary.BigEndian.Uint32(serialnumber_byte)

	result, _ := reader.ReadByte()

	return &FeedbackPacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnumber,
		Result:       result,
	}
}
