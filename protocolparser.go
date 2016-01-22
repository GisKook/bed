package bed

import (
	"bytes"
	"errors"
	"log"

	"github.com/giskook/bed/pb"
	"github.com/golang/protobuf/proto"
)

func CheckSum(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := uint16(1); i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}

func CheckProtocol(buffer *bytes.Buffer) (uint8, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if buffer.Bytes()[0] != 0xBA {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if bufferlen > 2 {
		pkglen := buffer.Bytes()[1]
		if pkglen < 7 { // flag + messagelen + cmdid + checksum + flag = 7  2048 is a magic number
			buffer.ReadByte()
			CheckProtocol(buffer)
		}
		if int(pkglen) > bufferlen {
			return HalfPack, 0
		} else {
			checksum := CheckSum(buffer.Bytes()[2:], uint16(pkglen))
			if checksum == buffer.Bytes()[bufferlen-2] && buffer.Bytes()[bufferlen-1] == 0xED {
				cmdid := buffer.Bytes()[2]
				return cmdid, uint16(pkglen+4) // 1 for ba 1 for len 1 for xor 1 for ed
			} else {
				buffer.ReadByte()
				CheckProtocol(buffer)
			}
		}
	} else {
		return HalfPack, 0
	}

	return Illegal, 0
}

func CheckNsqProtocol(message []byte) (uint64, uint32, *Report.Command, error) {
	command := &Report.ControlReport{}
	err := proto.Unmarshal(message, command)
	if err != nil {
		log.Println("unmarshal error")
		return 0, 0, nil, errors.New("unmarshal error")
	} else {
		bedid := command.Tid
		serialnum := command.SerialNumber
		cmd := command.GetCommand()

		return bedid, serialnum, cmd, nil
	}
}
