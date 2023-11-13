package model

type Instance struct {
	valid       bool
	instanceId  string
	port        uint32
	ip          string
	weight      uint32
	metadata    map[string]string
	clusterName string
	serviceName string
	enable      bool
	healthy     bool
	ephemeral   bool
}

func (i *Instance) GetIp() string {
	return i.ip
}

func (i *Instance) GetPort() uint32 {
	return i.port
}

func (i *Instance) GetWeight() uint32 {
	return i.weight
}

func (i *Instance) GetMetadata() map[string]string {
	return i.metadata
}

func (i *Instance) IsHealthy() bool {
	return i.healthy
}

func (i *Instance) IsEnabled() bool {
	return i.enable
}

func (i *Instance) GetClusterName() string {
	return i.clusterName
}

type Service struct {
	namespace       string
	group           string
	name            string
	ephemeral       bool
	revision        int64
	lastUpdatedTime int64
}

func (s *Service) GetGroup() string {
	return s.group
}

func (s *Service) GetName() string {
	return s.name
}

func (s *Service) GetNamespace() string {
	return s.namespace
}

func (s *Service) GetRevision() int64 {
	return s.revision
}
