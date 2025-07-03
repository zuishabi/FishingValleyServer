package routers

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"github.com/zuishabi/zinx/utils"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"go.uber.org/zap"
)

// OnUserReadyRouter
// 2
// 当客户端准备好时发出，调用注册好的函数
type OnUserReadyRouter struct {
	znet.BaseRouter
}

func (o *OnUserReadyRouter) Handle(request ziface.IRequest) {
	user := core.Omap.GetUserByConn(request.GetConnection())
	user.OnUserReady()
	// 将用户加入一个区域
	core.AddUserToMap(user, 1)
	if err := user.SendMsg(2, &FishingValleyProto.UserReady{}); err != nil {
		utils.L.Error("发送消息错误", zap.Error(err))
	}
}
