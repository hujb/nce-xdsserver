package bootstrap

import (
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nce/nce-xdsserver/common/event"
	"github.com/nce/nce-xdsserver/common/event/process"
	"github.com/nce/nce-xdsserver/common/resource"
	"github.com/nce/nce-xdsserver/nacos"
	"github.com/nce/nce-xdsserver/xds"
	"time"
)

func NewXdsServer() {

}

var sc *constant.ServerConfig

func InitXdsService() (*xds.NacosXdsService, *time.Ticker, chan *event.Event) {
	nacosXdsService := xds.NewNacosXdsService()
	eventProcess, eventsChan := process.NewEventProcessor(nacosXdsService)
	go eventProcess.HandleEvents()
	watcher, watcherTicker := nacos.NewNacosServiceInfoResourceWatcher(eventProcess)
	go watcher.ExecuteTimerTask()
	resourceManager := resource.GetResourceManagerInstance()
	resourceManager.SetResourceWatcher(watcher)
	nacosXdsService.SetResourceManager(resourceManager)
	return nacosXdsService, watcherTicker, eventsChan
}
