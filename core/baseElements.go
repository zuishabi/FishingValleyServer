package core

import (
	"sync"
	"sync/atomic"
)

// UserView 玩家的状态快照
type UserView struct {
	UID        uint32
	UserName   string
	Coordinate Coordinate
}

type Coordinate struct {
	X int32
	Y int32
}

type PlayerState struct {
	Coordinate Coordinate
	AreaID     uint32
	ChunkID    Coordinate
	PosLock    sync.RWMutex // 对玩家的位置和所处地图进行加锁
}

// SyncDirty 将脏数据同步
func (p *PlayerState) SyncDirty(dirty *DirtyPlayerState) {
	p.PosLock.Lock()
	defer p.PosLock.Unlock()
	if dirty.Need.Load() {
		p.Coordinate = dirty.coordinate.Load().(Coordinate)
	}
	dirty.Need.Store(false)
}

type DirtyPlayerState struct {
	coordinate atomic.Value //Coordinate
	Need       atomic.Bool
}

func (d *DirtyPlayerState) UpdateCoordinate(c Coordinate) {
	d.Need.Store(true)
	d.coordinate.Store(c)
}
