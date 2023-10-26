package common

import (
	"fmt"
	"gitlab2.psbc.com/ecn/ecn-xdsserver/common/resource"
	"gitlab2.psbc.com/ecn/ecn-xdsserver/xds"
	"time"
)

type EventProcessor struct {
	nacosXdsService *xds.NacosXdsService
	eventTicker     *time.Ticker
}

func NewEventProcessor(n *xds.NacosXdsService) (*EventProcessor, *time.Ticker) {
	ep := &EventProcessor{nacosXdsService: n}
	ep.eventTicker = time.NewTicker(3 * time.Second)
	return ep, ep.eventTicker
}

func (e *EventProcessor) ExecuteTimerTask() {
	time.Sleep(15 * time.Second)
	for {
		// 等待定时器触发事件
		<-e.eventTicker.C
		// 执行的定时任务
		rs := resource.GetInstance().CreateResourceSnapshot()
		e.nacosXdsService.HandChangedEvent(rs)
		fmt.Println("收到nacos服务变更通知，下发推送任务执行了！")
	}
}
