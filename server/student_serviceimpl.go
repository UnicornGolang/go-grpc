package main

// 首先我们需要导入服务协议
import (
	"context"
	"errors"
	"fmt"
	"ggrpc/idl/pb"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// 为了保证我们的服务安全性，我们需要对接口限流，防止服务瞬时请求数过高
// 巧妙的使用管道的容量的特性, 这里的容量就是我们希望的并大容量，因为超
// 过管道容量的时候开始阻塞，不会进入到方法实现，从而实现对数据库的保护
// 管道作为一个容器，我们让它传输一个空的结构体
var limit = make(chan struct{}, 10)

// 实现最终 idl 种定义的函数接口的结构体
type StudentServer struct{}

// 实现 idl 文件中定义的获取学生信息的服务
// 这里我们会在函数中获取到比定义的时候更多的参数，这个参数是
// grpc 框架为我们封装的请求的上下文对象，里面可以获取到一些请求的相信信息
func (s *StudentServer) GetStudentInfo(ctx context.Context, request *student_service.Request) (*student_service.Student, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("执行接口时出错: %v\n", err)
		}
	}()

	studentId := request.StudentId

	if len(studentId) == 0 {
		return nil, errors.New("studentId is empty")
	}
	fmt.Printf("接收到请求，请求参数: %v\n", studentId)

	// 获取接口处理的信息, 每次访问这个接口的时候，我们先写入一个空的结构体
	limit <- struct{}{}
	student := GetStudentInfo(studentId)
	// 在访问结束后从容量池里面取走一个值，这样就可以实现对应的接口并发性保护
	<-limit
	return &student, nil
}

// 获取学生信息
func GetStudentInfo(studentId string) student_service.Student {

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6380",
		Password: "redis123!@#",
		DB:       0,
	})

	ctx := context.TODO()
	student := student_service.Student{}
	for field, value := range client.HGetAll(ctx, studentId).Val() {
		fmt.Printf("%v-%v\n", field, value)
		switch field {
		case "Name":
			student.Name = value
		case "Age":
			age, err := strconv.Atoi(value)
			if err == nil {
				student.Age = int32(age)
			}
		case "Height":
			heigt, err := strconv.ParseFloat(value, 32)
			if err == nil {
				student.Height = float32(heigt)
			}
		}
	}
	return student
}
