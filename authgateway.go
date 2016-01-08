package sha

import (
	"fmt"
	"sync"
)

type BedProperty struct {
	Uid             uint64
	Bedversion      uint8
	Protocolversion uint8
}

type BedHub struct {
	Bed map[uint64]*BedProperty

	waitGroup *sync.WaitGroup
}

var bedhub *BedHub

func (g *BedHub) Add(bedid uint64, bedversion uint8, protocolversion uint8) {
	g.Bed[gatewayid] = &BedProperty{
		Uid:             bedid,
		Bedversion:      bedversion,
		Protocolversion: protocolversion,
	}
}

func NewBedHub() *BedHub {
	if bedhub == nil {
		bedhub = &BedHub{
			Bed:       make(map[uint64]*BedProperty),
			waitGroup: &sync.WaitGroup{},
		}
	}

	return bedhub
}

func (g *BedHub) Remove(bedid uint64) {
	delete(g.Bed, bedid)
}

func (g *BedHub) GetBed(bedid uint64) *BedProperty {
	bed, _ := g.Bed[bedid]

	return bed
}

func (g *BedHub) Check(bedid uint64) bool {
	_, ok := g.Bed[bedid]

	return ok
}
