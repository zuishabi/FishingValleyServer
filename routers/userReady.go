package routers

import (
	"FishingValleyServer/core"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
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
}
