package bed

import (
	"bytes"
	"github.com/giskook/gotcp"
	"log"
	"time"
)

var ConnSuccess uint8 = 0
var ConnUnauth uint8 = 1

type ConnConfig struct {
	HeartBeat    uint8
	ReadLimit    int64
	WriteLimit   int64
	NsqChanLimit int32
}

type Conn struct {
	conn                 *gotcp.Conn
	config               *ConnConfig
	recieveBuffer        *bytes.Buffer
	ticker               *time.Ticker
	readflag             int64
	writeflag            int64
	packetNsqReceiveChan chan gotcp.Packet
	closeChan            chan bool
	index                uint32
	uid                  uint64
	status               uint8
}

func NewConn(conn *gotcp.Conn, config *ConnConfig) *Conn {
	return &Conn{
		conn:                 conn,
		recieveBuffer:        bytes.NewBuffer([]byte{}),
		config:               config,
		readflag:             time.Now().Unix(),
		writeflag:            time.Now().Unix(),
		ticker:               time.NewTicker(time.Duration(config.HeartBeat) * time.Second * 10),
		packetNsqReceiveChan: make(chan gotcp.Packet, config.NsqChanLimit),
		closeChan:            make(chan bool),
		index:                0,
		status:               ConnSuccess,
	}
}

func (c *Conn) Close() {
	c.closeChan <- true
	c.ticker.Stop()
	c.recieveBuffer.Reset()
	close(c.packetNsqReceiveChan)
	close(c.closeChan)
}

func (c *Conn) GetBedID() uint64 {
	return c.uid
}
func (c *Conn) GetBuffer() *bytes.Buffer {
	return c.recieveBuffer
}

func (c *Conn) writeToclientLoop() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case p := <-c.packetNsqReceiveChan:
			if p != nil {
				c.conn.GetRawConn().Write(p.Serialize())
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) SendToBed(p gotcp.Packet) {
	c.packetNsqReceiveChan <- p
}

func (c *Conn) UpdateReadflag() {
	c.readflag = time.Now().Unix()
}

func (c *Conn) UpdateWriteflag() {
	c.writeflag = time.Now().Unix()
}

func (c *Conn) SetStatus(status uint8) {
	c.status = status
}

func (c *Conn) checkHeart() {
	defer func() {
		c.conn.Close()
	}()

	var now int64
	for {
		select {
		case <-c.ticker.C:
			now = time.Now().Unix()
			if now-c.readflag > c.config.ReadLimit {
				log.Println("read linmit")
				return
			}
			if now-c.writeflag > c.config.WriteLimit {
				log.Println("write limit")
				return
			}
			if c.status == ConnUnauth {
				log.Printf("close connection %d\n", c.uid)
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) Do() {
	go c.checkHeart()
	go c.writeToclientLoop()
}

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	heartbeat := GetConfiguration().GetServerConnCheckInterval()
	readlimit := GetConfiguration().GetServerReadLimit()
	writelimit := GetConfiguration().GetServerWriteLimit()
	config := &ConnConfig{
		HeartBeat:  uint8(heartbeat),
		ReadLimit:  int64(readlimit),
		WriteLimit: int64(writelimit),
	}
	conn := NewConn(c, config)

	c.PutExtraData(conn)

	NewConns().Add(conn)
	conn.Do()

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData().(*Conn)
	conn.Close()
	NewConns().Remove(conn.GetBedID())
	NewBedHub().Remove(conn.GetBedID())
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	bedPacket := p.(*BedPacket)
	switch bedPacket.Type {
	case Login:
		c.AsyncWritePacket(bedPacket, time.Second)
	case HeartBeat:
		c.AsyncWritePacket(bedPacket, time.Second)
	case AppControlFeedback:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case HandleControlFeedback:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case AppPottyFeedback:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case HandlePottyFeedback:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case AfterPotty:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	case AppBedReset:
		GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	}

	return true
}
