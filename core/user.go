package core

import (
	"github.com/zuishabi/zinx/utils"
	"github.com/zuishabi/zinx/ziface"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"sync"
)

type User struct {
	UID         uint32
	UserName    string
	DirtyState  DirtyPlayerState
	PlayerState PlayerState
	Conn        ziface.IConnection
	sync.Mutex
}

func (u *User) SendMsg(msgID uint32, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return u.Conn.SendBuffMsg(msgID, data)
}

func (u *User) Move(x int32, y int32, direction bool) {
	u.DirtyState.UpdateCoordinate(Coordinate{
		X: x,
		Y: y,
	}, direction)
}

// SendMsgAround 向周围玩家发送消息
func (u *User) SendMsgAround(msgID uint32, msg proto.Message) {
	u.PlayerState.PosLock.RLock()
	defer u.PlayerState.PosLock.RUnlock()
	area := Areas[u.PlayerState.AreaID]
	if area == nil {
		return
	}
	for i := range area.Grids[u.PlayerState.ChunkID.X][u.PlayerState.ChunkID.Y].users {
		if i != u.UID {
			if err := area.Grids[u.PlayerState.ChunkID.X][u.PlayerState.ChunkID.Y].users[i].SendMsg(msgID, msg); err != nil {
				utils.L.Error("发送消息失败", zap.Error(err))
			}
		}
	}
	if u.PlayerState.ChunkID.X-1 >= 0 {
		for i := range area.Grids[u.PlayerState.ChunkID.X-1][u.PlayerState.ChunkID.Y].users {
			if err := area.Grids[u.PlayerState.ChunkID.X-1][u.PlayerState.ChunkID.Y].users[i].SendMsg(msgID, msg); err != nil {
				utils.L.Error("发送消息失败", zap.Error(err))
			}
		}
		if u.PlayerState.ChunkID.Y-1 >= 0 {
			for i := range area.Grids[u.PlayerState.ChunkID.X-1][u.PlayerState.ChunkID.Y-1].users {
				if err := area.Grids[u.PlayerState.ChunkID.X-1][u.PlayerState.ChunkID.Y-1].users[i].SendMsg(msgID, msg); err != nil {
					utils.L.Error("发送消息失败", zap.Error(err))
				}
			}
		}
		if u.PlayerState.ChunkID.Y+1 < area.LGridNum {
			for i := range area.Grids[u.PlayerState.ChunkID.X-1][u.PlayerState.ChunkID.Y+1].users {
				if err := area.Grids[u.PlayerState.ChunkID.X-1][u.PlayerState.ChunkID.Y+1].users[i].SendMsg(msgID, msg); err != nil {
					utils.L.Error("发送消息失败", zap.Error(err))
				}
			}
		}
	}
	if u.PlayerState.ChunkID.X+1 < area.WGridNum {
		for i := range area.Grids[u.PlayerState.ChunkID.X+1][u.PlayerState.ChunkID.Y].users {
			if err := area.Grids[u.PlayerState.ChunkID.X+1][u.PlayerState.ChunkID.Y].users[i].SendMsg(msgID, msg); err != nil {
				utils.L.Error("发送消息失败", zap.Error(err))
			}
		}
		if u.PlayerState.ChunkID.Y-1 >= 0 {
			for i := range area.Grids[u.PlayerState.ChunkID.X+1][u.PlayerState.ChunkID.Y-1].users {
				if err := area.Grids[u.PlayerState.ChunkID.X+1][u.PlayerState.ChunkID.Y-1].users[i].SendMsg(msgID, msg); err != nil {
					utils.L.Error("发送消息失败", zap.Error(err))
				}
			}
		}
		if u.PlayerState.ChunkID.Y+1 < area.LGridNum {
			for i := range area.Grids[u.PlayerState.ChunkID.X+1][u.PlayerState.ChunkID.Y+1].users {
				if err := area.Grids[u.PlayerState.ChunkID.X+1][u.PlayerState.ChunkID.Y+1].users[i].SendMsg(msgID, msg); err != nil {
					utils.L.Error("发送消息失败", zap.Error(err))
				}
			}
		}
	}
	if u.PlayerState.ChunkID.Y-1 >= 0 {
		for i := range area.Grids[u.PlayerState.ChunkID.X][u.PlayerState.ChunkID.Y-1].users {
			if err := area.Grids[u.PlayerState.ChunkID.X][u.PlayerState.ChunkID.Y-1].users[i].SendMsg(msgID, msg); err != nil {
				utils.L.Error("发送消息失败", zap.Error(err))
			}
		}
	}
	if u.PlayerState.ChunkID.Y+1 < area.LGridNum {
		for i := range area.Grids[u.PlayerState.ChunkID.X][u.PlayerState.ChunkID.Y+1].users {
			if err := area.Grids[u.PlayerState.ChunkID.X][u.PlayerState.ChunkID.Y+1].users[i].SendMsg(msgID, msg); err != nil {
				utils.L.Error("发送消息失败", zap.Error(err))
			}
		}
	}
}

var readyFunctionList []func(user *User)

func AddReadyFunction(f func(user *User)) {
	readyFunctionList = append(readyFunctionList, f)
}

func (u *User) OnUserReady() {
	for _, i := range readyFunctionList {
		i(u)
	}
}
