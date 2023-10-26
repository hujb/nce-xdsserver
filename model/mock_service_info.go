package model

type MockServiceInfo struct {
	hosts []*Instance
}

func (m *MockServiceInfo) GetHosts() []*Instance {
	return m.hosts
}
