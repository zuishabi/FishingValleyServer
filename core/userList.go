package core

import (
	"github.com/zuishabi/zinx/ziface"
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
