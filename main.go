package main

import (
	"flag"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/nce/nce-xdsserver/bootstrap"
	"github.com/nce/nce-xdsserver/common/event"
	"github.com/nce/nce-xdsserver/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

var (
	grpcAddr = flag.String("grpcAddr", "127.0.0.1:28848", "Address of the xds server")
	dc       = flag.String("dc", "Y", "The center of data")
	workId   = flag.String("workId", "99999", "the instance num in system")
)

func init() {
	flag.Parse()
	//日志配置初始化
	log.InitLogger(*dc, *workId)
}

func main() {
	nacosXdsService, watcherTicker, eventsChan := bootstrap.InitXdsService()
	server := grpc.NewServer()
	listen, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		closeResource(err, watcherTicker, eventsChan)
		panic(err)
	}
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(server, nacosXdsService)
	//log.Printf("Starting XdsServer, grpcAddr=%v", *grpcAddr)
	log.Logger.Info("MSG=XdsServer listening, grpcAddr=" + *grpcAddr)
	err = server.Serve(listen)
	if err != nil {
		closeResource(err, watcherTicker, eventsChan)
		os.Exit(-1)
	}
}

func closeResource(err error, watcherTicker *time.Ticker, eventsChan chan *event.Event) {
	//log.Printf("XdsServer started failed, err: %v", err)
	log.Logger.Error("MSG=XdsServer started failed, err:" + err.Error())

	defer watcherTicker.Stop()
	defer close(eventsChan)
}
