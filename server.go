package main

import (
	"FishingValleyServer/core"
	FishingValleyProto "FishingValleyServer/protobuf"
	"FishingValleyServer/routers"
	"fmt"
	ZSC "github.com/zuishabi/ServiceCenter/src"
	"github.com/zuishabi/zinx/utils"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"time"
)

func onConnStop(connection ziface.IConnection) {
	user := core.Omap.GetUserByConn(connection)
	if user == nil {
		fmt.Println("当前用户不存在")
		return
	}
	// 广播退出信息
	core.Omap.RemoveUser(user.UID)
}

// 122.228.237.118:34896
func main() {
	go gameLoop()
	core.InitAreaMap(1, 320, 10, 5, core.Coordinate{
		X: 118,
		Y: 120,
	})
	c := ZSC.NewClient("127.0.0.1", 8880, "FishingValley", nil, "127.0.0.1:9999")
	if err := c.Start(); err != nil {
		fmt.Println(err)
		return
	}
	s := znet.NewServer()
	s.SetOnConnStop(onConnStop)
	utils.InitLogger(0)
	s.AddRouter(1, &routers.LoginRouter{})
	s.AddRouter(3, &routers.MovementRouter{})
	s.AddRouter(2, &routers.OnUserReadyRouter{})
	s.AddRouter(5, &routers.PlayerChunkChangeRouter{})
	s.AddRouter(10, &routers.PlayerNameRouter{})
	s.Serve()
}

// 游戏主循环
func gameLoop() {
	for {
		start := time.Now()
		users := core.Omap.GetAllUsers()
		// 将玩家脏数据同步
		for _, v := range users {
			if v != nil {
				v.PlayerState.SyncDirty(&v.DirtyState)
			}
		}
		// 同步玩家坐标
		posMsg := FishingValleyProto.Movement{}
		for _, v := range users {
			if v != nil {
				v.PlayerState.PosLock.RLock()
				posMsg.X = v.PlayerState.Coordinate.X
				posMsg.Y = v.PlayerState.Coordinate.Y
				posMsg.Uid = v.UID
				v.PlayerState.PosLock.RUnlock()
				v.SendMsgAround(3, &posMsg)
			}
		}
		duration := time.Since(start)
		if duration < time.Millisecond*50 {
			time.Sleep(time.Millisecond*50 - duration)
		}
	}
}
