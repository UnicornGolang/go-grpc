package main

import (
	student_service "ggrpc/idl/pb"
	"net"

	"google.golang.org/grpc"
)

// 测试服务端代码的编写
func main() {

	// 使用 net 包创建一个端口监听
	listen, err := net.Listen("tcp", ":2346")
	if err != nil {
		panic(err)
	}

	// 使用 grpc 框架创建服务
	server := grpc.NewServer()

	// 将我们的实现的服务接口注册到服务列表中
	// RegisterStudentServiceServer 这个注解的函数是 idl 文件中为
	// student_service 定义的函数生成的注册器
	student_service.RegisterStudentServiceServer(server, new(StudentServer))

	// 将监听与服务进行关联
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
