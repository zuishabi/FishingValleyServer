package core

import (
	"fmt"
	"github.com/zuishabi/zinx/ziface"
	"google.golang.org/protobuf/proto"
	"sync"
)

type OnlineMap struct {
	sync.RWMutex
	um map[uint32]*User
}

var Omap *OnlineMap

func init() {
	Omap = &OnlineMap{
		um: make(map[uint32]*User),
	}
}

func (m *OnlineMap) AddUser(user *User) {
	m.Lock()
	defer m.Unlock()
	m.um[user.UID] = user
}

func (m *OnlineMap) RemoveUser(uid uint32) {
	m.Lock()
	defer m.Unlock()
	delete(m.um, uid)
	fmt.Println("移除用户")
}

func (m *OnlineMap) GetUser(uid uint32) *User {
	m.RLock()
	defer m.RUnlock()
	return m.um[uid]
}

func (m *OnlineMap) GetUserByConn(conn ziface.IConnection) *User {
	m.RLock()
	defer m.RUnlock()
	uid, err := conn.GetProperty("uid")
	if err != nil {
		return nil
	}
	return m.um[uid.(uint32)]
}

// GetAllUsersView 返回用户的状态快照
func (m *OnlineMap) GetAllUsersView() (users []UserView) {
	m.RLock()
	defer m.RUnlock()
	users = make([]UserView, 0)
	for _, user := range m.um {
		users = append(users, UserView{
			UID:        user.UID,
			UserName:   user.UserName,
			Coordinate: user.coordinate,
		})
	}
	return
}

// GetAllUsers 获得所有用户
func (m *OnlineMap) GetAllUsers() (users []*User) {
	m.RLock()
	defer m.RUnlock()
	users = make([]*User, 0)
	for _, user := range m.um {
		users = append(users, user)
	}
	return
}

func (m *OnlineMap) BroadCast(msgID uint32, message proto.Message) {
	m.RLock()
	defer m.RUnlock()
	for _, user := range m.um {
		if err := user.SendMsg(msgID, message); err != nil {
			continue
		}
	}
}
