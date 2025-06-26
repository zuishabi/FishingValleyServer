package routers

import (
	"FishingValleyServer/core"
	LoginCenterProto2 "FishingValleyServer/loginCenterGRPC"
	FishingValleyProto "FishingValleyServer/protobuf"
	"context"
	"fmt"
	"github.com/zuishabi/zinx/ziface"
	"github.com/zuishabi/zinx/znet"
	"google.golang.org/protobuf/proto"
)

// LoginRouter
// 1
// 当用户登录后验证其登陆码是否一致，然后将其添加到user列表中
type LoginRouter struct {
	znet.BaseRouter
}

func (l *LoginRouter) Handle(request ziface.IRequest) {
	loginMsg := &FishingValleyProto.ConfirmLogin{}
	_ = proto.Unmarshal(request.GetData(), loginMsg)
	//从客户端获得的key中在redis中查找是否存在
	//s := database.RDB.Get(context.Background(), loginMsg.Key)
	//id, err := database.Client.Do(context.Background(), database.Client.B().Get().Key(loginMsg.Key).Build()).AsInt64()
	rsp, err := LoginCenterProto2.LoginCenterCheckLoginClient.CheckLogin(context.Background(), &LoginCenterProto2.CheckLoginReq{Key: loginMsg.Key})
	if err != nil {
		sendLoginFailMsg(request.GetConnection(), "登录失败，凭证失效")
		return
	}
	//查找当前用户是否在线
	if core.Omap.GetUser(rsp.Uid) != nil {
		sendLoginFailMsg(request.GetConnection(), "当前用户已登录")
		return
	}
	//创建用户到在线列表
	iUser := core.User{
		UID:      rsp.Uid,
		Conn:     request.GetConnection(),
		UserName: rsp.UserName,
	}
	request.GetConnection().SetProperty("uid", rsp.Uid)
	core.Omap.AddUser(&iUser)
	sendLoginSuccessMsg(request.GetConnection(), rsp.Uid)
}

// 向客户端发送登录失败的消息
func sendLoginFailMsg(conn ziface.IConnection, errMsg string) {
	loginMsg := &FishingValleyProto.ConfirmLoginResponse{
		Success: false,
		Content: errMsg,
	}
	m, _ := proto.Marshal(loginMsg)
	err := conn.SendBuffMsg(1, m)
	if err != nil {
		fmt.Println("[LoginRouter sendLoginFailMsg] : SendMsg err = ", err)
	}
}

// 向客户端发送登录成功的消息
func sendLoginSuccessMsg(conn ziface.IConnection, uid uint32) {
	loginMsg := &FishingValleyProto.ConfirmLoginResponse{
		Success: true,
		Id:      uid,
	}
	m, err := proto.Marshal(loginMsg)
	if err != nil {
		fmt.Println("[LoginRouter sendLoginSuccessMsg] : proto marshal err = ", err)
		return
	}
	err = conn.SendBuffMsg(1, m)
	if err != nil {
		fmt.Println("[LoginRouter sendLoginSuccessMsg] : SendMsg err = ", err)
	}
}
