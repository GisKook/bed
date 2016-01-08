package sha

import (
	"encoding/binary"
	"github.com/giskook/smarthome-access/pb"
)

type NsqSetDevicename struct {
	GatewayID    uint64
	SerialNumber uint32
	DeviceID     uint64
	Name         string
}

func (p *NsqSetDevicename) Serialize() []byte {
	buf := []byte{0xCE, 0x00, 0x00, 0x80, 0x08}
	gatewayid := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid, p.GatewayID)
	buf = append(buf, gatewayid[2:]...)
	serialnum := make([]byte, 4)
	binary.BigEndian.PutUint32(serialnum, p.SerialNumber)
	buf = append(buf, serialnum...)
	buf = append(buf, 0x06)
	deviceid := make([]byte, 8)
	binary.BigEndian.PutUint64(deviceid, p.DeviceID)
	buf = append(buf, deviceid[2:]...)
	buf = append(buf, byte(len(p.Name)))
	buf = append(buf, []byte(p.Name)...)
	totallen := len(buf) + 2 // 2 for chechsum and end flag
	binary.BigEndian.PutUint16(buf[1:3], uint16(totallen))
	buf = append(buf, CheckSum(buf, uint16(len(buf))))
	buf = append(buf, 0xCE)

	return buf
}

func ParseNsqSetDevicename(gatewayid uint64, serialnum uint32, command *Report.Command) *NsqSetDevicename {
	cmdparam := command.GetParas()
	return &NsqSetDevicename{
		GatewayID:    gatewayid,
		SerialNumber: serialnum,
		DeviceID:     cmdparam[0].Npara,
		Name:         cmdparam[1].Strpara,
	}
}
