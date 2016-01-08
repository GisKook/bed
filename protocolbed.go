package sha

import (
	"github.com/giskook/gotcp"
)

var (
	Illegal  uint16 = 0
	HalfPack uint16 = 255

	Login                 uint16 = 1
	HeartBeat             uint16 = 2
	AppControlFeedback    uint16 = 3
	HandleControlFeedback uint16 = 4
	AppPottyFeedback      uint16 = 5
	HandlePottyFeedback   uint16 = 6
	AfterPotty            uint16 = 7
)

type BedPacket struct {
	Type   uint16
	Packet gotcp.Packet
}

func (this *BedPacket) Serialize() []byte {
	switch this.Type {
	case Login:
		return this.Packet.(*LoginPacket).Serialize()
	case HeartBeat:
		return this.Packet.(*HeartPacket).Serialize()
	case AppControlFeedback:
		return this.Packet.(*FeedbackAppControlPacket).Serialize()
	case HandleControlFeedback:
		return this.Packet.(*FeedbackHandleControlPacket).Serialize()
	case AppPottyFeedback:
		return this.Packet.(*FeedbackAppPottyPacket).Serialize()
	case HandlePottyFeedback:
		return this.Packet.(*FeedbackeHandlePottyPacket).Serialize()
	case AfterPotty:
		return this.Packet.(*AfterPottyPacket).Serialize()
	}

	return nil
}

func NewBedPacket(Type uint16, Packet gotcp.Packet) *BedPacket {
	return &BedPacket{
		Type:   Type,
		Packet: Packet,
	}
}

type BedProtocol struct {
}

func (this *BedProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	smconn := c.GetExtraData().(*Conn)
	smconn.UpdateReadflag()

	buffer := smconn.GetBuffer()

	conn := c.GetRawConn()
	for {
		data := make([]byte, 2048)
		readLengh, err := conn.Read(data)

		if err != nil {
			return nil, err
		}

		if readLengh == 0 {
			return nil, gotcp.ErrConnClosing
		} else {
			buffer.Write(data[0:readLengh])
			cmdid, pkglen := CheckProtocol(buffer)
			//		log.Printf("recv box cmd %d \n", cmdid)

			pkgbyte := make([]byte, pkglen)
			buffer.Read(pkgbyte)
			switch cmdid {
			case Login:
				pkg := ParseLogin(pkgbyte, smconn)
				return NewBedPacket(Login, pkg), nil
			case HeartBeat:
				pkg := ParseHeart(pkgbyte)
				return NewBedPacket(HeartBeat, pkg), nil
			case AppControlFeedback:
				pkg := ParseAppControl(pkgbyte, smconn)
				return NewBedPacket(AppControlFeedback, pkg), nil
			case HandleControlFeedback:
				pkg := ParseHandleControlFeedback(pkgbyte)
				return NewBedPacket(HandleControlFeedback, pkg), nil
			case AppPottyFeedback:
				pkg := ParseAppPottyFeedback(pkgbyte)
				return NewBedPacket(AppPottyFeedback, pkg), nil
			case HandlePottyFeedback:
				pkg := ParseHandlePottyFeedback(pkgbyte)
				return NewBedPacket(HandlePottyFeedback, pkg), nil
			case AfterPotty:
				pkg := ParseAfterPotty(pkgbyte)
				return NewBedPacket(AfterPotty, pkg), nil

			case Illegal:
			case HalfPack:
			}
		}
	}

}
