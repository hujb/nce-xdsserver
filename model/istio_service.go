package model

import (
	"github.com/nce/nce-xdsserver/nacos/nacosResource"
	"time"
)

type IstioService struct {
	Name            string
	GroupName       string
	Namespace       string
	Revision        int64
	Hosts           []*nacosResource.NacosInstance
	CreateTimeStamp time.Time
}

//func NewIstioService(service *Service, mockServiceInfo *MockServiceInfo) *IstioService {
//	return &IstioService{
//		name:            service.GetName(),
//		groupName:       service.GetGroup(),
//		namespace:       service.GetNamespace(),
//		revision:        service.GetRevision(),
//		createTimeStamp: time.Now(),
//		hosts:           sanitizeServiceInfo(mockServiceInfo),
//	}
//}
//
//func NewIstioServiceByOld(service *Service, mockServiceInfo *MockServiceInfo, old *IstioService) *IstioService {
//	return &IstioService{
//		name:            service.GetName(),
//		groupName:       service.GetGroup(),
//		namespace:       service.GetNamespace(),
//		revision:        service.GetRevision(),
//		createTimeStamp: old.GetCreateTimeStamp(),
//		hosts:           sanitizeServiceInfo(mockServiceInfo),
//	}
//}

//func sanitizeServiceInfo(mockServiceInfo *MockServiceInfo) []*nacosResource.NacosInstance {
//	hosts := make([]*nacosResource.NacosInstance, 0)
//	// TODO
//	return hosts
//}

//func (m *IstioService) GetCreateTimeStamp() time.Time {
//	return m.CreateTimeStamp
//}
//
//func (m *IstioService) GetHosts() []*nacosResource.NacosInstance {
//	return m.Hosts
//}
//
//func (m *IstioService) GetNamespace() string {
//	return m.Namespace
//}
//
//func (m *IstioService) GetRevision() int64 {
//	return m.Revision
//}
