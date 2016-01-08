package sha

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

type NsqCheckOnlinePacket struct {
	GatewayID    uint64
	SerialNumber uint32
	Result       uint8
}

func (p *NsqCheckOnlinePacket) Serialize() []byte {
	rep := &Report.ControlReport{
		Tid:          p.GatewayID,
		SerialNumber: p.SerialNumber,
		Command: &Report.Command{
			Type: Report.Command_CMT_REPONLINE,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(p.Result),
				},
			},
		},
	}

	repdata, err := proto.Marshal(rep)

	if err != nil {
		log.Println("marshaling error", err)
	}

	return repdata
}

func ParseNsqCheckOnline(gatewayid uint64, serialnum uint32) *NsqCheckOnlinePacket {
	var result uint8 = Offline
	if NewConns().Check(gatewayid) {
		result = Online
	}

	return &NsqCheckOnlinePacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
		Result:       result,
	}
}
