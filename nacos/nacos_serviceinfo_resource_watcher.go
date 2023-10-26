package nacos

import (
	"fmt"
	"gitlab2.psbc.com/ecn/ecn-xdsserver/model"
	"time"
)

type NacosServiceInfoResourceWatcher struct {
	serviceInfoMap map[string]*model.IstioService
	callbacks      []func()
	watcherTicker  *time.Ticker
}

func NewNacosServiceInfoResourceWatcher() (*NacosServiceInfoResourceWatcher, *time.Ticker) {
	watcher := &NacosServiceInfoResourceWatcher{
		serviceInfoMap: map[string]*model.IstioService{},
		callbacks:      []func(){},
		watcherTicker:  time.NewTicker(3 * time.Second),
	}

	watcher.SubscribeAllServices(func() {
		// TODO  Query all services to see if any of them have changes.

		// TODO 变更判断，更新serviceInfoMap

		// TODO 变更通知
	})

	// TODO channel通知事件队列在这里做？

	return watcher, watcher.watcherTicker
}

func (w *NacosServiceInfoResourceWatcher) ExecuteTimerTask() {
	time.Sleep(15 * time.Second)
	for {
		// 等待定时器触发事件
		<-w.watcherTicker.C
		// 定时订阅nacos服务信息
		for _, callback := range w.callbacks {
			callback()
		}

		fmt.Println("定时查询Nacos服务信息任务执行了！")
	}
}

func (w *NacosServiceInfoResourceWatcher) SubscribeAllServices(SubscribeCallback func()) {
	w.callbacks = append(w.callbacks, SubscribeCallback)
}

// Snapshot TODO 优化？
func (w *NacosServiceInfoResourceWatcher) Snapshot() map[string]*model.IstioService {
	var cloneMap = make(map[string]*model.IstioService)
	for key, value := range w.serviceInfoMap {
		cloneMap[key] = value
	}
	return cloneMap
}
