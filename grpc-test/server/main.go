package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	pb "mypb.com/pb"
	"net"
)

// hello server

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello," + in.Name}, nil
}

func (s *server) SayManyHello(stream pb.Greeter_SayManyHelloServer) error {
	for {
		// 接收流式请求
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// 将收到的字符串反转并返回
		reverse := func(s string) string {
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}
		// 如果接受到"shutdown",停止接受流请求
		if in.GetName() == "shutdown" {
			return nil
		}

		// 创建响应并返回
		reply := &pb.HelloReply{
			Message: reverse(in.GetName()),
		}
		if err := stream.Send(reply); err != nil {
			return err
		}

	}
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
