package main

import (
	"context"
	pb "github.com/hujb/nce-xdsserver/helloworld/proto/grpc/service"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:16848", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	res, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "我是客户端,祝贺你"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("服务端发来的信息：%s", res.GetMessage())
}
