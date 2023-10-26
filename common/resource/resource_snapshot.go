package resource

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 需要保证是私有的，不允许被直接访问更改，故需要通过对象访问，那么是都通过一个对象访问还是允许每个线程创一个对象，这里的考量点是提高并发性能，故允许每个线程创一个对象，
// 那么当主线程需要的时候，必须要保证拿到最新的版本号，那么需要提供一个用于获取保存了当前最新版本号的对象，这个需要引入一个全局单例的是资源管理器ResourceManager，提供ResourceSnapshot对象的获取和创建
var versionSuffix = 0

type ResourceSnapshot struct {
	version     string
	isCompleted bool
	//使用sync.RWMutex可以允许多个goroutine同时读取共享资源，提高并发性能，但在写入时需要独占访问，确保数据的一致性和安全性。在不需要读取并发访问的情况下，可以使用sync.Mutex简化代码和性能开销
	mutex sync.RWMutex
}

func (r *ResourceSnapshot) InitResourceSnapshot() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isCompleted {
		return
	}
	r.generateVersion()
	r.isCompleted = true
}

func (r *ResourceSnapshot) generateVersion() {
	timeNow := fmt.Sprintf("%v", time.Now())
	versionSuffix++
	r.version = timeNow + "/" + strconv.Itoa(versionSuffix)
}

func (r *ResourceSnapshot) GetVersion() string {
	return r.version
}
