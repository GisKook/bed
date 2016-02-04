package bed

import (
	"errors"
	"github.com/giskook/gotcp"
	"log"
)

var (
	Illegal  uint8 = 254
	HalfPack uint8 = 253

	Login                 uint8 = 0
	HeartBeat             uint8 = 255
	AppControlFeedback    uint8 = 1
	HandleControlFeedback uint8 = 2
	AppPottyFeedback      uint8 = 3
	HandlePottyFeedback   uint8 = 4
	AfterPotty            uint8 = 5
	AppBedReset           uint8 = 6

	ErrNotLogin = errors.New("do not login")
)

type BedPacket struct {
	Type   uint8
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
		return this.Packet.(*FeedbackAppControlPacket).Serialize()
	case AppPottyFeedback:
		return this.Packet.(*FeedbackPottyPacket).Serialize()
	case HandlePottyFeedback:
		return this.Packet.(*FeedbackPottyPacket).Serialize()
	case AfterPotty:
		return this.Packet.(*FeedbackAfterPottyPacket).Serialize()
	case AppBedReset:
		return this.Packet.(*FeedbackAppControlPacket).Serialize()
	}

	return nil
}

func NewBedPacket(Type uint8, Packet gotcp.Packet) *BedPacket {
	return &BedPacket{
		Type:   Type,
		Packet: Packet,
	}
}

type BedProtocol struct {
}

func (this *BedProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	smconn := c.GetExtraData().(*Conn)

	buffer := smconn.GetBuffer()

	conn := c.GetRawConn()
	for {
		data := make([]byte, 2048)
		readLengh, err := conn.Read(data)
		log.Printf("recv %x", data[0:readLengh])

		if err != nil {
			return nil, err
		}

		if readLengh == 0 {
			return nil, gotcp.ErrConnClosing
		} else {
			buffer.Write(data[0:readLengh])
			cmdid, pkglen := CheckProtocol(buffer)
			log.Println(cmdid)
			log.Println(NewConns().Check(0))
			if cmdid != Login && !NewConns().Check(smconn.GetBedID()) {
				return nil, ErrNotLogin
			}

			smconn.UpdateReadflag()

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
				pkg := ParseAppControlFeedback(pkgbyte, smconn, AppControlFeedback)
				return NewBedPacket(AppControlFeedback, pkg), nil
			case HandleControlFeedback:
				bedpkg := ParseHandleControlFeedback(pkgbyte)
				smconn.SendToBed(bedpkg)
				pkg := ParseAppControlFeedback(pkgbyte, smconn, HandleControlFeedback)
				return NewBedPacket(HandleControlFeedback, pkg), nil
			case AppPottyFeedback:
				pkg := ParsePottyFeedback(pkgbyte, smconn, AppPottyFeedback)
				return NewBedPacket(AppPottyFeedback, pkg), nil
			case HandlePottyFeedback:
				bedpkg := ParseHandlePottyFeedback(pkgbyte)
				smconn.SendToBed(bedpkg)
				pkg := ParsePottyFeedback(pkgbyte, smconn, HandlePottyFeedback)
				return NewBedPacket(HandlePottyFeedback, pkg), nil
			case AfterPotty:
				bedpkg := ParseAfterPottyTobedFeedback(pkgbyte)
				smconn.SendToBed(bedpkg)
				pkg := ParseAfterPottyFeedback(pkgbyte, smconn)
				return NewBedPacket(AfterPotty, pkg), nil
			case AppBedReset:
				pkg := ParseAppControlFeedback(pkgbyte, smconn, AppBedReset)
				return NewBedPacket(AppBedReset, pkg), nil

			case Illegal:
			case HalfPack:
			}
		}
	}

}
