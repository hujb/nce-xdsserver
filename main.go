package main

import (
	"fmt"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/hujb/nce-xdsserver/common"
	"github.com/hujb/nce-xdsserver/common/resource"
	"github.com/hujb/nce-xdsserver/nacos"
	"github.com/hujb/nce-xdsserver/xds"
	"google.golang.org/grpc"
	"net"
)

func init() {
	//TODO
	fmt.Println("执行进程初始化函数...")
}

func main() {
	nacosXdsService := initXdsService()
	server := grpc.NewServer()
	listen, err := net.Listen("tcp", "192.168.0.102:28848")
	if err != nil {
		panic(err)
	}
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(server, nacosXdsService)
	watcher, watcherTicker := nacos.NewNacosServiceInfoResourceWatcher()
	go watcher.ExecuteTimerTask()
	eventProcess, eventTicker := common.NewEventProcessor(nacosXdsService)
	go eventProcess.ExecuteTimerTask()
	err = server.Serve(listen)
	if err != nil {
		watcherTicker.Stop()
		eventTicker.Stop()
		return
	}
}

func initXdsService() *xds.NacosXdsService {
	nacosXdsService := xds.NewNacosXdsService()
	resource.GetInstance()
	return nacosXdsService
}
