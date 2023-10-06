package main

import (
	"context"
	"fmt"
	student_service "ggrpc/idl/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 进行服务的调用
	for i := 0; i < 100; i++ {
		callService()
	}
}

func callService() {
	conn, err := grpc.Dial("127.0.0.1:2346", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("connect to server error: %v\n", err)
		return
	}
	defer conn.Close()
	client := student_service.NewStudentServiceClient(conn)

	resp, err := client.GetStudentInfo(context.TODO(), &student_service.Request{StudentId: "SG-2029120212202"})
	if err != nil {
		fmt.Printf("调用 grpc 服务失败: %v\n", err)
	}
	fmt.Printf("{ Name : '%s', Age : %d, Height : %.1f}\n", resp.Name, resp.Age, resp.Height)
}
