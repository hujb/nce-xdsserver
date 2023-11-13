package resource

import (
	"github.com/nce/nce-xdsserver/model"
	"github.com/nce/nce-xdsserver/nacos"
	"sync"
)

// ResourceManager ResourceSnapshot由资源管理器提供访问和创建
type ResourceManager struct {
	resourceSnapshot           *ResourceSnapshot
	mutex                      sync.RWMutex
	serviceInfoResourceWatcher *nacos.NacosServiceInfoResourceWatcher
}

var singletonResourceManager *ResourceManager
var once sync.Once

// GetResourceManagerInstance sync.Once类型，它提供了一种保证函数只会被执行一次的机制。在GetInstance()函数中，我们使用了once.Do()方法，并将创建单例对象的逻辑放在了括号中的函数中。这样，在每次调用GetInstance()时，都会检查once的值，如果为Once{}，则执行括号中的函数创建单例对象，并将其赋值给instance，然后返回instance。如果once的值不为Once{}，则直接返回instance。
func GetResourceManagerInstance() *ResourceManager {
	once.Do(func() {
		singletonResourceManager = &ResourceManager{resourceSnapshot: &ResourceSnapshot{isCompleted: false}}
	})
	return singletonResourceManager
}

func (rmg *ResourceManager) SetResourceWatcher(watcher *nacos.NacosServiceInfoResourceWatcher) {
	rmg.serviceInfoResourceWatcher = watcher
}

func (rmg *ResourceManager) Services() map[string]*model.IstioService {
	return rmg.serviceInfoResourceWatcher.Snapshot()
}

func (rmg *ResourceManager) GetResourceSnapshot() *ResourceSnapshot {
	rmg.mutex.Lock()
	defer rmg.mutex.Unlock()
	return rmg.resourceSnapshot
}

func (rmg *ResourceManager) SetResourceSnapshot(rs *ResourceSnapshot) {
	rmg.mutex.Lock()
	defer rmg.mutex.Unlock()
	rmg.resourceSnapshot = rs
}

// InitResourceSnapshot 主线程启动时调用
func (rmg *ResourceManager) InitResourceSnapshot() {
	rs := rmg.GetResourceSnapshot()
	rs.InitResourceSnapshot(rmg)
}

// CreateResourceSnapshot 每个协程启动时调用
func (rmg *ResourceManager) CreateResourceSnapshot() *ResourceSnapshot {
	rs := &ResourceSnapshot{isCompleted: false}
	rs.InitResourceSnapshot(rmg)
	rmg.SetResourceSnapshot(rs)
	return rs
}
