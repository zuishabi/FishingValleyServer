package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"fmt"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type PlayerChunkChangeRouter struct {
	znet.BaseRouter
}

func (p *PlayerChunkChangeRouter) Handle(request ziface.IRequest) {
	msg := FishingValleyProto.PlayerChunkChange{}
	_ = proto.Unmarshal(request.GetData(), &msg)
	// 首先获得当前玩家的区域
	user := core.Omap.GetUserByConn(request.GetConnection())
	user.PlayerState.PosLock.Lock()
	areaID := user.PlayerState.AreaID
	area := core.Areas[areaID]
	area.Grids[user.PlayerState.ChunkID.X][user.PlayerState.ChunkID.Y].RemoveUser(user)
	user.PlayerState.ChunkID = core.Coordinate{
		X: int32(msg.ChunkX),
		Y: int32(msg.ChunkY),
	}
	user.PlayerState.PosLock.Unlock()
	fmt.Println("将玩家", user.UID, "添加到区块", msg.ChunkX, ",", msg.ChunkY)
	area.Grids[msg.ChunkX][msg.ChunkY].AddUser(user)
}
