package main

import (
	"flag"
	"fmt"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/nce/nce-xdsserver/bootstrap"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	grpcAddr = flag.String("grpcAddr", "127.0.0.1:28848", "Address of the xds server")
)

func init() {
	//TODO
	fmt.Println("执行进程初始化函数...")
}

func main() {
	flag.Parse()
	nacosXdsService, watcherTicker, eventsChan := bootstrap.InitXdsService()
	server := grpc.NewServer()
	listen, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		panic(err)
	}
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(server, nacosXdsService)
	err = server.Serve(listen)
	if err != nil {
		defer watcherTicker.Stop()
		defer close(eventsChan)
		os.Exit(-1)
	}
}
