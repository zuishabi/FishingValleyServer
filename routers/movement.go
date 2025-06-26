package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type MovementRouter struct {
	znet.BaseRouter
}

func (m *MovementRouter) Handle(request ziface.IRequest) {
	user := core.Omap.GetUserByConn(request.GetConnection())
	movement := FishingValleyProto.Movement{}
	_ = proto.Unmarshal(request.GetData(), &movement)
	user.Move(movement.X, movement.Y)
}
