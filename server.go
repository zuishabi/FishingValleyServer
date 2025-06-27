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
	core.Omap.RemoveUser(user.UID)
}

// 122.228.237.118:34896
func main() {
	go syncUsers()
	c := ZSC.NewClient("122.228.237.118", 34896, "FishingValley", nil, "127.0.0.1:9999")
	if err := c.Start(); err != nil {
		fmt.Println(err)
		return
	}
	s := znet.NewServer()
	s.SetOnConnStop(onConnStop)
	utils.InitLogger(0)
	s.AddRouter(1, &routers.LoginRouter{})
	s.AddRouter(3, &routers.MovementRouter{})
	s.AddRouter(4, &routers.PlayerNameRouter{})
	s.Serve()
}

func syncUsers() {
	for {
		select {
		case <-time.After(time.Millisecond * 50):
			users := core.Omap.GetAllUsersView()
			coordinates := make([]FishingValleyProto.Movement, len(users))
			for i, v := range users {
				coordinates[i] = FishingValleyProto.Movement{
					Uid: v.UID,
					X:   v.Coordinate.X,
					Y:   v.Coordinate.Y,
				}
			}
			for _, v := range core.Omap.GetAllUsers() {
				if v != nil {
					for i, _ := range coordinates {
						if v.UID != users[i].UID {
							if err := v.SendMsg(3, &coordinates[i]); err != nil {
								fmt.Println(err)
							}
						}
					}
				}
			}
		}
	}
}
