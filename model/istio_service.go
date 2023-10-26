package model

import "time"

type IstioService struct {
	name            string
	groupName       string
	namespace       string
	revision        int64
	hosts           []*Instance
	createTimeStamp time.Time
}

func NewIstioService(service *Service, mockServiceInfo *MockServiceInfo) *IstioService {
	return &IstioService{
		name:            service.GetName(),
		groupName:       service.GetGroup(),
		namespace:       service.GetNamespace(),
		revision:        service.GetRevision(),
		createTimeStamp: time.Now(),
		hosts:           sanitizeServiceInfo(mockServiceInfo),
	}
}

func NewIstioServiceByOld(service *Service, mockServiceInfo *MockServiceInfo, old *IstioService) *IstioService {
	return &IstioService{
		name:            service.GetName(),
		groupName:       service.GetGroup(),
		namespace:       service.GetNamespace(),
		revision:        service.GetRevision(),
		createTimeStamp: old.GetCreateTimeStamp(),
		hosts:           sanitizeServiceInfo(mockServiceInfo),
	}
}

func sanitizeServiceInfo(mockServiceInfo *MockServiceInfo) []*Instance {
	hosts := make([]*Instance, 0)
	// TODO
	return hosts
}

func (m *IstioService) GetCreateTimeStamp() time.Time {
	return m.createTimeStamp
}
