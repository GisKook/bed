package sha

import (
	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

var (
	ModifySuceess uint8 = 0
	OldPasswdErr  uint8 = 1
	ModifyFail    uint8 = 2
)

type NsqChangePasswdPacket struct {
	GatewayID    uint64
	SerialNumber uint32
	OldPasswd    string
	NewPasswd    string
	Result       uint8
}

func (p *NsqChangePasswdPacket) Serialize() []byte {
	para := []*Report.Command_Param{
		&Report.Command_Param{
			Type:  Report.Command_Param_UINT8,
			Npara: uint64(p.Result),
		},
	}

	command := &Report.Command{
		Type:  Report.Command_CMT_REQCHANGEPASSWD,
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

func ParseNsqChangePasswd(gatewayid uint64, serialnum uint32, command *Report.Command) *NsqChangePasswdPacket {
	commandparam := command.GetParas()
	log.Println("change passwd ")

	var result uint8 = ModifySuceess
	if GetUserPasswdHub().Check(gatewayid, commandparam[0].Strpara) {
		log.Println("change passwd %s %s\n", commandparam[0].Strpara, commandparam[1].Strpara)
		err := GetServer().GetDatabase().SetPasswd("passwd", gatewayid, commandparam[1].Strpara)
		if err != nil {
			result = ModifyFail
			log.Println(err.Error())
		}
	} else {
		result = OldPasswdErr
	}

	return &NsqChangePasswdPacket{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
		OldPasswd:    commandparam[0].Strpara,
		NewPasswd:    commandparam[1].Strpara,
		Result:       result,
	}
}
