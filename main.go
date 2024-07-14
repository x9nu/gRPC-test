// main 文件只是测试一下序列化和反序列化
// 真正使用的是 `grpc-server.go` 和 `grpc-client.go“
package main

import (
	"fmt"
	"grpc-test/service"

	"google.golang.org/protobuf/proto"
)

func main() {
	user := &service.User{
		Username: "test",
		Age:      18,
	}
	// 序列化的过程
	marshal, err := proto.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化：", marshal)

	// 反序列化的过程
	userData := &service.User{}
	err = proto.Unmarshal(marshal, userData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("反序列化：", userData)

}
