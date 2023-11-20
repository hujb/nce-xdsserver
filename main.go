package main

import (
	"flag"
	"fmt"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/nce/nce-xdsserver/bootstrap"
	"github.com/nce/nce-xdsserver/common/event"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
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
		closeResource(err, watcherTicker, eventsChan)
		panic(err)
	}
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(server, nacosXdsService)
	log.Printf("Starting XdsServer, grpcAddr=%v", *grpcAddr)
	err = server.Serve(listen)
	if err != nil {
		closeResource(err, watcherTicker, eventsChan)
		os.Exit(-1)
	}
}

func closeResource(err error, watcherTicker *time.Ticker, eventsChan chan *event.Event) {
	log.Printf("XdsServer started failed, err: %v", err)
	defer watcherTicker.Stop()
	defer close(eventsChan)
}
