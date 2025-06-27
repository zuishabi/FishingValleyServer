package core

// UserView 玩家的状态快照
type UserView struct {
	UID        uint32
	UserName   string
	Coordinate Coordinate
}

type Coordinate struct {
	X int32
	Y int32
}
