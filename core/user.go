package core

import (
	"github.com/zuishabi/zinx/ziface"
	"google.golang.org/protobuf/proto"
	"sync"
)

type User struct {
	UID        uint32
	UserName   string
	coordinate Coordinate
	Conn       ziface.IConnection
	sync.Mutex
}

func (u *User) SendMsg(msgID uint32, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return u.Conn.SendBuffMsg(msgID, data)
}

func (u *User) Move(x int32, y int32) {
	u.coordinate.X = x
	u.coordinate.Y = y
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
