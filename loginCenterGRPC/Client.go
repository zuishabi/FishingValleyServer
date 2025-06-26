package LoginCenterProto

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var opts []grpc.DialOption = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
var LoginCenterCheckLoginClient CheckLoginServiceClient
var LoginCenterUserClient UserServiceClient

func init() {
	c, err := grpc.NewClient(fmt.Sprintf("%s:%d", "127.0.0.1", 7777), opts...)
	if err != nil {
		//连接失败
		panic(err)
	}
	LoginCenterCheckLoginClient = NewCheckLoginServiceClient(c)
	LoginCenterUserClient = NewUserServiceClient(c)
}
