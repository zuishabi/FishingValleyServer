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

/*
	状态表:
	0、普通状态
	1、等待钓鱼状态
	2、鱼上钩状态
	3、成功钓鱼状态
*/

type PlayerState struct {
	Coordinate   Coordinate
	AreaID       uint32
	ChunkID      Coordinate
	Direction    bool         // 玩家面对的方向
	PosLock      sync.RWMutex // 对玩家的位置和所处地图进行加锁
	CurrentState uint32       // 玩家当前状态
	StateLock    sync.RWMutex // 对玩家当前状态进行加锁
}

// SyncDirty 将脏数据同步
func (p *PlayerState) SyncDirty(user *User, dirty *DirtyPlayerState) {
	p.PosLock.Lock()
	defer p.PosLock.Unlock()
	ok := true
	p.Direction = dirty.direction.Load()
	p.Coordinate, ok = dirty.coordinate.Load().(Coordinate)
	if !ok {
		return
	}
	chunk, ok := dirty.chunkID.Load().(Coordinate)
	if !ok {
		return
	}
	if p.ChunkID != chunk {
		areaID := p.AreaID
		area := Areas[areaID]
		area.Grids[p.ChunkID.X][p.ChunkID.Y].RemoveUser(user)
		p.ChunkID = chunk
		area.Grids[chunk.X][chunk.Y].AddUser(user)
	}
}

type DirtyPlayerState struct {
	coordinate atomic.Value //Coordinate
	direction  atomic.Bool
	chunkID    atomic.Value //Coordinate
}

func (d *DirtyPlayerState) UpdateCoordinate(c Coordinate, direction bool) {
	d.coordinate.Store(c)
	d.direction.Store(direction)
}

func (d *DirtyPlayerState) UpdateChunkID(chunkID Coordinate) {
	d.chunkID.Store(chunkID)
}
