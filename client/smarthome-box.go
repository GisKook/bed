package shb

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"
)

func CheckSum(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := uint16(1); i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}

type Device struct {
	DeviceID   uint64
	DeviceType uint8
	Company    uint16
	Name       string
	Status     uint8
}

type Smarthomebox struct {
	GatewayID   uint64
	Name        string
	DeviceCount uint16
	DeviceList  []*Device
	Wg          *sync.WaitGroup
	ExitChan    chan struct{}
}

func char2byte(c string) byte {
	switch c {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "a":
		return 10
	case "b":
		return 11
	case "c":
		return 12
	case "d":
		return 13
	case "e":
		return 14
	case "f":
		return 15
	}
	return 0
}

func Macaddr2uint64(mac string) uint64 {
	var buffer []byte
	buffer = append(buffer, 0)
	buffer = append(buffer, 0)
	value := char2byte(string(mac[0]))*16 + char2byte(string(mac[1]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[2]))*16 + char2byte(string(mac[3]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[4]))*16 + char2byte(string(mac[5]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[6]))*16 + char2byte(string(mac[7]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[8]))*16 + char2byte(string(mac[9]))
	buffer = append(buffer, value)
	value = char2byte(string(mac[10]))*16 + char2byte(string(mac[11]))
	buffer = append(buffer, value)

	return binary.BigEndian.Uint64(buffer)
}
func NewSmarthomebox(strgatewayid string, name string) *Smarthomebox {
	gatewayid := Macaddr2uint64(strgatewayid)

	return &Smarthomebox{
		GatewayID:   gatewayid,
		Name:        name,
		DeviceCount: 0,
		DeviceList:  nil,
		Wg:          &sync.WaitGroup{},
		ExitChan:    make(chan struct{}),
	}
}

func (b *Smarthomebox) Close() {
	close(b.ExitChan)
}

func (b *Smarthomebox) Check(deviceid uint64) bool {
	for i := uint16(0); i < b.DeviceCount; i++ {
		if b.DeviceList[i].DeviceID == deviceid {
			return true
		}
	}

	return false
}

func (b *Smarthomebox) Del(deviceid uint64) {
	for i := uint16(0); i < b.DeviceCount; i++ {
		if b.DeviceList[i].DeviceID == deviceid {
			b.DeviceCount--
			b.DeviceList[i] = nil
			b.DeviceList = append(b.DeviceList[:i], b.DeviceList[i+1:]...)
			return
		}
	}
}

func (b *Smarthomebox) Add(strdeviceid string, devicetype uint8, company uint16, name string, status uint8) {
	deviceid := Macaddr2uint64(strdeviceid)
	device := &Device{
		DeviceID:   deviceid,
		DeviceType: devicetype,
		Company:    company,
		Name:       name,
		Status:     status,
	}
	if b.Check(deviceid) {
		b.Del(deviceid)
		b.DeviceCount--
	}

	b.DeviceList = append(b.DeviceList, device)
	b.DeviceCount++
}

func (b *Smarthomebox) login(conn *net.TCPConn) uint8 {
	logincmd := []byte{0xCE, 0x00, 0x00, 0x00, 0x01}
	gatewayid_byte := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid_byte, b.GatewayID)
	logincmd = append(logincmd, gatewayid_byte[2:]...)
	logincmd = append(logincmd, byte(len(b.Name)))
	logincmd = append(logincmd, []byte(b.Name)...)
	logincmd = append(logincmd, byte(1))
	logincmd = append(logincmd, byte(1))
	devicecount_byte := make([]byte, 2)
	binary.BigEndian.PutUint16(devicecount_byte, b.DeviceCount)
	logincmd = append(logincmd, devicecount_byte...)
	for i := uint16(0); i < b.DeviceCount; i++ {
		logincmd = append(logincmd, byte(b.DeviceList[i].DeviceType))
		logincmd = append(logincmd, byte(6))
		deviceid_byte := make([]byte, 8)
		binary.BigEndian.PutUint64(deviceid_byte, b.DeviceList[i].DeviceID)
		logincmd = append(logincmd, deviceid_byte[2:]...)
		logincmd = append(logincmd, 0x00)
		logincmd = append(logincmd, 0x01)
		logincmd = append(logincmd, 0x01)
		logincmd = append(logincmd, byte(len(b.DeviceList[i].Name)))
		logincmd = append(logincmd, []byte(b.DeviceList[i].Name)...)
	}
	cmdlen := len(logincmd) + 2 // 2 for checksum and end flag
	binary.BigEndian.PutUint16(logincmd[1:3], uint16(cmdlen))
	logincmd = append(logincmd, CheckSum(logincmd, uint16(cmdlen-2)))
	logincmd = append(logincmd, 0xCE)

	log.Printf("%X\n", logincmd)
	_, err := conn.Write(logincmd)
	if err != nil {
		log.Println(err.Error())

		return 0
	}

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	if buffer[11] == 0x01 {
		log.Println("Login success")
		return 1
	} else {
		log.Println("Login fail")
		return 0
	}
}

func (b *Smarthomebox) heart(conn *net.TCPConn) {
	b.Wg.Add(1)
	defer func() {
		b.Wg.Done()
	}()
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-b.ExitChan:
			log.Println("heart exit")
			return
		case <-ticker.C:
			heartcmd := []byte{0xCE, 0x00, 0x0D, 0x00, 0x02}
			heart_byte := make([]byte, 8)
			binary.BigEndian.PutUint64(heart_byte, b.GatewayID)
			heartcmd = append(heartcmd, heart_byte[2:]...)
			heartcmd = append(heartcmd, CheckSum(heartcmd, uint16(len(heartcmd))))
			heartcmd = append(heartcmd, 0xCE)
			log.Printf("%X\n", heartcmd)
			_, err := conn.Write(heartcmd)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

}

func (b *Smarthomebox) adddeldevice(conn *net.TCPConn) {
	b.Wg.Add(1)
	defer func() {
		b.Wg.Done()
	}()
	ticker := time.NewTicker(3 * time.Second)
	add := true
	for {
		select {
		case <-b.ExitChan:
			log.Println("adddel exit")
			return
		case <-ticker.C:
			adddelcmd := []byte{0xCE, 0x00, 0x00, 0x00, 0x05}
			gatewayid_byte := make([]byte, 8)
			binary.BigEndian.PutUint64(gatewayid_byte, b.GatewayID)
			adddelcmd = append(adddelcmd, gatewayid_byte[2:]...)
			if add {
				adddelcmd = append(adddelcmd, byte(1))
				add = false
			} else {
				adddelcmd = append(adddelcmd, byte(0))
				add = true
			}
			adddelcmd = append(adddelcmd, byte(1))
			adddelcmd = append(adddelcmd, byte(6))
			adddelcmd = append(adddelcmd, []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0xEE, 0x00, 0x01, 0x01}...)
			cmdlen := len(adddelcmd) + 2 // 2 for checksum and end flag
			binary.BigEndian.PutUint16(adddelcmd[1:3], uint16(cmdlen))
			adddelcmd = append(adddelcmd, CheckSum(adddelcmd, uint16(cmdlen-2)))
			adddelcmd = append(adddelcmd, 0xCE)
			log.Printf("add or del %X\n", adddelcmd)
			_, err := conn.Write(adddelcmd)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (b *Smarthomebox) warnup(conn *net.TCPConn) {
	b.Wg.Add(1)
	defer func() {
		b.Wg.Done()
	}()

	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-b.ExitChan:
			log.Println("adddel exit")
			return
		case <-ticker.C:
			adddelcmd := []byte{0xCE, 0x00, 0x1E, 0x00, 0x06}
			gatewayid_byte := make([]byte, 8)
			binary.BigEndian.PutUint64(gatewayid_byte, b.GatewayID)
			adddelcmd = append(adddelcmd, gatewayid_byte[2:]...)
			adddelcmd = append(adddelcmd, byte(1))
			adddelcmd = append(adddelcmd, byte(6))
			adddelcmd = append(adddelcmd, []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0xEE}...)
			t := time.Now().Unix()
			t_byte := make([]byte, 8)
			binary.BigEndian.PutUint64(t_byte, uint64(t))
			adddelcmd = append(adddelcmd, t_byte...)
			adddelcmd = append(adddelcmd, byte(1))
			cmdlen := len(adddelcmd) + 2 // 2 for checksum and end flag
			binary.BigEndian.PutUint16(adddelcmd[1:3], uint16(cmdlen))
			adddelcmd = append(adddelcmd, CheckSum(adddelcmd, uint16(cmdlen-2)))
			adddelcmd = append(adddelcmd, 0xCE)
			log.Printf("warn up %X\n", adddelcmd)
			_, err := conn.Write(adddelcmd)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

}

func (b *Smarthomebox) setname(deviceid uint64, name string) {
	for i := uint16(0); i < b.DeviceCount; i++ {
		if b.DeviceList[i].DeviceID == deviceid {
			b.DeviceList[i].Name = name
		}
	}
}

func (b *Smarthomebox) setnamefeedback(conn *net.TCPConn, buffer []byte) {
	snfb := []byte{0xCE, 0x00, 0x00, 0x00, 0x08}
	gatewayid_byte := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid_byte, b.GatewayID)
	snfb = append(snfb, gatewayid_byte[2:]...)
	snfb = append(snfb, buffer[11:15]...)
	snfb = append(snfb, 0x01)
	snfb = append(snfb, 0x06)
	snfb = append(snfb, buffer[16:22]...)
	namelen := buffer[22]
	snfb = append(snfb, buffer[22])
	snfb = append(snfb, buffer[23:23+namelen]...)
	deviceid_byte := []byte{0x00, 0x00}
	deviceid_byte = append(deviceid_byte, buffer[16:22]...)
	deviceid := binary.BigEndian.Uint64(deviceid_byte)
	b.setname(deviceid, string(buffer[23:23+namelen]))
	cmdlen := len(snfb) + 2 // 2 for checksum and end flag
	binary.BigEndian.PutUint16(snfb[1:3], uint16(cmdlen))
	snfb = append(snfb, CheckSum(snfb, uint16(cmdlen-2)))
	snfb = append(snfb, 0xCE)
	log.Printf("set name %X\n", snfb)
	_, err := conn.Write(snfb)
	if err != nil {
		log.Println(err.Error())
	}

}

func (b *Smarthomebox) opfeedback(conn *net.TCPConn, buffer []byte) {
	opfb := []byte{0xCE, 0x00, 0x00, 0x00, 0x04}
	gatewayid_byte := make([]byte, 8)
	binary.BigEndian.PutUint64(gatewayid_byte, b.GatewayID)
	opfb = append(opfb, gatewayid_byte[2:]...)
	log.Printf("op serial %X\n", buffer[0:100])

	opfb = append(opfb, buffer[11:15]...)
	opfb = append(opfb, 0x01)
	cmdlen := len(opfb) + 2 // 2 for checksum and end flag
	binary.BigEndian.PutUint16(opfb[1:3], uint16(cmdlen))
	opfb = append(opfb, CheckSum(opfb, uint16(cmdlen-2)))
	opfb = append(opfb, 0xCE)

	log.Printf("op feedback %X\n", opfb)
	_, err := conn.Write(opfb)
	if err != nil {
		log.Println(err.Error())
	}

}

func (b *Smarthomebox) recv(conn *net.TCPConn) {
	b.Wg.Add(1)
	defer func() {
		b.Wg.Done()
	}()
	for {
		select {
		case <-b.ExitChan:
			return
		default:
		}

		buffer := make([]byte, 1024)
		length, _ := conn.Read(buffer)
		if buffer[3] == 0x80 && buffer[4] == 0x08 {
			b.setnamefeedback(conn, buffer)
		} else if buffer[3] == 0x80 && buffer[4] == 0x04 {
			b.opfeedback(conn, buffer)
		}
		log.Printf("recv %X\n", buffer[0:length])
	}
}

func (b *Smarthomebox) Do(srvaddr string, wg *sync.WaitGroup) {
	wg.Add(1)
	b.Wg.Add(1)

	tcpaddr, _ := net.ResolveTCPAddr("tcp", srvaddr)

	conn, err := net.DialTCP("tcp", nil, tcpaddr)

	defer func() {
		b.Wg.Done()
		conn.Close()
		wg.Done()
	}()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if b.login(conn) == 1 {
		go b.heart(conn)
		go b.recv(conn)
		go b.adddeldevice(conn)
		go b.warnup(conn)
	}
	b.Wg.Wait()
}
