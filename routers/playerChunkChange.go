package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
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
	user.DirtyState.UpdateChunkID(core.Coordinate{
		X: int32(msg.ChunkX),
		Y: int32(msg.ChunkY),
	})
}
