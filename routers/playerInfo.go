package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"fmt"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type PlayerInfoRouter struct {
	znet.BaseRouter
}

func (p *PlayerInfoRouter) Handle(request ziface.IRequest) {
	req := FishingValleyProto.PlayerInfoReq{}
	_ = proto.Unmarshal(request.GetData(), &req)
	user := core.Omap.GetUserByConn(request.GetConnection())
	u := core.Omap.GetUser(req.Uid)
	u.PlayerState.StateLock.RLock()
	defer u.PlayerState.StateLock.RUnlock()
	rsp := FishingValleyProto.PlayerInfoRsp{
		Uid:       u.UID,
		Name:      u.UserName,
		State:     u.PlayerState.CurrentState,
		Direction: u.PlayerState.Direction,
	}
	if err := user.SendMsg(10, &rsp); err != nil {
		fmt.Println(err)
	}
}
