package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"fmt"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type PlayerNameRouter struct {
	znet.BaseRouter
}

func (p *PlayerNameRouter) Handle(request ziface.IRequest) {
	req := FishingValleyProto.PlayerNameReq{}
	_ = proto.Unmarshal(request.GetData(), &req)
	user := core.Omap.GetUserByConn(request.GetConnection())
	u := core.Omap.GetUser(req.Uid)
	fmt.Println("req.uid = ", req.Uid)
	rsp := FishingValleyProto.PlayerNameRsp{
		Uid:  u.UID,
		Name: u.UserName,
	}
	if err := user.SendMsg(4, &rsp); err != nil {
		fmt.Println(err)
	}
}
