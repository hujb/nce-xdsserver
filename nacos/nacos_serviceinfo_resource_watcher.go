package nacos

import (
	"flag"
	"fmt"
	"github.com/nce/nce-xdsserver/common/adapter"
	"github.com/nce/nce-xdsserver/common/event"
	"github.com/nce/nce-xdsserver/model"
	nceModel "github.com/nce/nce-xdsserver/model"
	"github.com/nce/nce-xdsserver/nacos/nacosResource"
	"github.com/nce/nce-xdsserver/util"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"sync"
	"time"
)

var (
	NacosAddr = flag.String("nacosAddr", "127.0.0.1:8848", "Address of the nacos server")
)

type NacosServiceInfoResourceWatcher struct {
	//serviceInfoMap map[string]*model.IstioService
	serviceInfoMap sync.Map
	callbacks      []func()
	watcherTicker  *time.Ticker
	eventAdapter   adapter.ResourceWatcherAdapter
}

func NewNacosServiceInfoResourceWatcher(eventAdapter adapter.ResourceWatcherAdapter) (*NacosServiceInfoResourceWatcher, *time.Ticker) {
	watcher := &NacosServiceInfoResourceWatcher{
		serviceInfoMap: sync.Map{},
		callbacks:      []func(){},
		watcherTicker:  time.NewTicker(3 * time.Second),
		eventAdapter:   eventAdapter,
	}
	flag.Parse()
	watcher.SubscribeAllServices(func() {
		var changed = false
		// 查询所有服务实例，是否有变更
		namespaces, err := GetAllNamespaces(*NacosAddr)
		if err != nil {
			return
		}
		for _, namespace := range namespaces {
			if namespace == "" {
				namespace = "public"
			}
			param := &nceModel.GetAllServiceInfoParam{NameSpace: namespace}
			serviceClusterInstanceData, err := GetAllServicesWithInstanceByNamespace(*NacosAddr, param)
			if err != nil {
				return
			}
			instances := make([]*nacosResource.NacosInstance, 0, 10)
			for _, serviceClusterInstanceDetail := range serviceClusterInstanceData {
				var serviceName string
				if len(serviceClusterInstanceDetail.ClusterMap) == 0 {
					continue
				}
				for clusterName, clusterMapDetail := range serviceClusterInstanceDetail.ClusterMap {
					if clusterMapDetail.Hosts == nil {
						continue
					}
					serviceName = util.BuildServiceNameForServiceEntry(serviceClusterInstanceDetail, clusterName, namespace)
					if serviceName == "istio-test10.FUN-A-MDP.DEFAULT-GROUP.public" {
						changed = true
					}
					for _, host := range clusterMapDetail.Hosts {
						host.ClusterName = clusterName
						host.ServiceName = serviceClusterInstanceDetail.ServiceName
						instances = append(instances, host)

					}
				}
				istioService := &model.IstioService{Name: serviceClusterInstanceDetail.ServiceName,
					Namespace: namespace,
					GroupName: serviceClusterInstanceDetail.GroupName,
					Hosts:     instances,
				}
				watcher.serviceInfoMap.Store(serviceName, istioService)
			}
		}

		// TODO 变更判断，更新serviceInfoMap

		// 变更通知
		if changed {
			watcher.eventAdapter.Notify(event.SERVICE_UPDATE_EVENT)
		}
	})

	return watcher, watcher.watcherTicker
}

func (w *NacosServiceInfoResourceWatcher) ExecuteTimerTask() {
	time.Sleep(5 * time.Second)
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

func (w *NacosServiceInfoResourceWatcher) Snapshot() map[string]*model.IstioService {
	cloneMap := make(map[string]*model.IstioService)
	w.serviceInfoMap.Range(func(key, value interface{}) bool {
		value = value.(*model.IstioService)
		data, err := json.Marshal(value)
		if err != nil {
			log.Fatalf("深度拷贝序列化异常，value: %v, err: %v", value, err)
			return false
		}
		cloneValue := &model.IstioService{}
		err = json.Unmarshal(data, cloneValue)
		if err != nil {
			log.Fatalf("深度拷贝反序列化异常，value: %v, err: %v", value, err)
			return false
		}
		cloneMap[key.(string)] = cloneValue
		return true
	})
	return cloneMap
}
