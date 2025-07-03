package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type SpeakRouter struct {
	znet.BaseRouter
}

func (s *SpeakRouter) Handle(request ziface.IRequest) {
	speak := FishingValleyProto.Speak{}
	_ = proto.Unmarshal(request.GetData(), &speak)
	user := core.Omap.GetUserByConn(request.GetConnection())
	speak.UserName = user.UserName
	core.Omap.BroadcastAll(7, &speak)
}
