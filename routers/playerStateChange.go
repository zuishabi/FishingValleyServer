package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type PlayerStateChangeRouter struct {
	znet.BaseRouter
}

func (p *PlayerStateChangeRouter) Handle(request ziface.IRequest) {
	user := core.Omap.GetUserByConn(request.GetConnection())
	from := FishingValleyProto.PlayerStateChange{}
	_ = proto.Unmarshal(request.GetData(), &from)
	from.Uid = user.UID
	// 将玩家动作进行转发
	user.SendMsgAround(8, &from)
	// 更新玩家状态
	user.PlayerState.StateLock.Lock()
	user.PlayerState.CurrentState = from.Action
	user.PlayerState.StateLock.Unlock()
}
