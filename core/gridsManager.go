package core

import (
	FishingValleyProto "FishingValleyServer/protobuf"
	"github.com/zuishabi/zinx/utils"
	"go.uber.org/zap"
	"sync"
)

// Grid 地图格
type Grid struct {
	StartPoint Coordinate
	users      map[uint32]*User
	sync.RWMutex
}

// AddUser 在格子中添加一个玩家
func (g *Grid) AddUser(user *User) {
	g.Lock()
	defer g.Unlock()
	// 将用户加入当前格子中
	g.users[user.UID] = user
}

// RemoveUser 在格子中删除一个玩家
func (g *Grid) RemoveUser(user *User) {
	g.Lock()
	defer g.Unlock()
	delete(g.users, user.UID)
}

// GetUsers 获得当前格子中所有的玩家
func (g *Grid) GetUsers() []*User {
	g.RLock()
	defer g.RUnlock()
	res := make([]*User, len(g.users))
	size := 0
	for i := range g.users {
		res[size] = g.users[i]
		size += 1
	}
	return res
}

// Areas 所有的区域列表
var Areas map[uint32]*AreaMap = make(map[uint32]*AreaMap)

// AreaMap 区域地图
type AreaMap struct {
	Size       uint32 // 单元格的大小
	Grids      [][]*Grid
	ID         uint32
	SpawnPoint Coordinate // 区域的出生点
	LGridNum   int32
	WGridNum   int32
}

func (a *AreaMap) GetGridFromCoordinate(c Coordinate) (int32, int32) {
	x := float32(c.X) / 10.0
	y := float32(c.Y) / 10.0
	return int32(uint32(x / float32(a.Size))), int32(y / float32(a.Size))
}

// InitAreaMap 生成地图区域
func InitAreaMap(id uint32, size uint32, lGridNum int32, wGridNum int32, spawnPoints Coordinate) {
	grids := make([][]*Grid, lGridNum)
	for i := range grids {
		grids[i] = make([]*Grid, wGridNum)
		for v := range grids[i] {
			grids[i][v] = &Grid{
				users: make(map[uint32]*User),
			}
		}
	}
	Areas[id] = &AreaMap{
		Size:       size,
		Grids:      grids,
		ID:         id,
		SpawnPoint: spawnPoints,
		LGridNum:   lGridNum,
		WGridNum:   wGridNum,
	}
}

// AddUserToMap 将玩家添加到一个新的区域
func AddUserToMap(user *User, id uint32) {
	user.PlayerState.PosLock.Lock()
	target := Areas[id]
	if user.PlayerState.AreaID != 0 {
		// 从一个区域转移到另一个区域，将用户从原先区域中取出
		source := Areas[user.PlayerState.AreaID]
		sourceGrid := source.Grids[user.PlayerState.ChunkID.X][user.PlayerState.ChunkID.Y]
		sourceGrid.RemoveUser(user)
	}
	tx, ty := target.GetGridFromCoordinate(target.SpawnPoint)
	targetGrid := target.Grids[tx][ty]
	targetGrid.AddUser(user)
	// 改变用户的区域id和坐标
	user.PlayerState.Coordinate = target.SpawnPoint
	user.PlayerState.AreaID = target.ID
	user.PlayerState.ChunkID = Coordinate{
		X: tx,
		Y: ty,
	}
	msg := FishingValleyProto.TransmitPlayer{
		AreaId: target.ID,
		X:      target.SpawnPoint.X * 10,
		Y:      target.SpawnPoint.Y * 10,
	}
	if err := user.SendMsg(6, &msg); err != nil {
		utils.L.Error("发送消息失败", zap.Error(err))
	}
	user.PlayerState.PosLock.Unlock()
}

// RemoveUserFromMap 将玩家从一个区域中移除
func RemoveUserFromMap(user *User) {
	user.PlayerState.PosLock.RLock()
	defer user.PlayerState.PosLock.RUnlock()
	if user.PlayerState.AreaID != 0 {
		source := Areas[user.PlayerState.AreaID]
		if source != nil {
			sourceGrid := source.Grids[user.PlayerState.ChunkID.X][user.PlayerState.ChunkID.Y]
			sourceGrid.RemoveUser(user)
		}
	}
	msg := FishingValleyProto.PlayerLeave{Uid: user.UID}
	user.SendMsgAround(4, &msg)
}
