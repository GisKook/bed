package sha

type ControlPacket struct {
	GatewayID uint64
	DeviceID  uint64
	Operate   uint8
}

func (p *ControlPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0x00)
	buf = append(buf, 0x00)
	buf = append(buf, 0x00)
	buf = append(buf, 0x00)

	return buf
}

func NewControlPacket(GatewayID uint64, DeviceID uint64, Operate uint8) *ControlPacket {
	return &ControlPacket{
		GatewayID: GatewayID,
		DeviceID:  DeviceID,
		Operate:   Operate,
	}
}
