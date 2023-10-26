package model

type Instance struct {
	instanceId  string
	port        uint64
	ip          string
	weight      float64
	metadata    map[string]string
	clusterName string
	serviceName string
	enable      bool
	healthy     bool
	ephemeral   bool
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
