package process

import (
	"github.com/nce/nce-xdsserver/common/event"
	"github.com/nce/nce-xdsserver/common/resource"
	"github.com/nce/nce-xdsserver/log"
	"github.com/nce/nce-xdsserver/xds"
)

type EventProcessor struct {
	nacosXdsService *xds.NacosXdsService
	//eventTicker     *time.Ticker
	eventsChan chan *event.Event
}

func NewEventProcessor(n *xds.NacosXdsService) (*EventProcessor, chan *event.Event) {
	ep := &EventProcessor{nacosXdsService: n}
	//ep.eventTicker = time.NewTicker(3 * time.Second)
	// 为了保证chan是阻塞的，没有给缓存容量
	ep.eventsChan = make(chan *event.Event)
	return ep, ep.eventsChan
}

func (e *EventProcessor) HandleEvents() {
	//time.Sleep(15 * time.Second) 2023.11.06废弃
	for {
		// 等待定时器触发事件 2023.11.06废弃
		//<-e.eventTicker.C
		select {
		case <-e.eventsChan:
			rs := resource.GetResourceManagerInstance().CreateResourceSnapshot()
			//e.nacosXdsService.HandChangedEvent(rs)
			e.nacosXdsService.HandleEvent(rs)
			//fmt.Println("收到nacos服务变更通知，下发推送任务执行完成！")
			log.Logger.Info("收到nacos服务变更通知，下发推送任务执行完成！")
		}
	}
}

func (e *EventProcessor) Notify(event *event.Event) {
	e.eventsChan <- event
}
