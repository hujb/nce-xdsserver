package model

type QueryAllServiceInfoParam struct {
	NamespaceId   string
	GroupName     string
	WithInstances bool
	PageNo        uint32
	PageSize      uint32
}

type QueryAllInstanceInfoByServiceParam struct {
	ServiceName string
	GroupName   string
	NamespaceId string
	Clusters    string
	HealthyOnly bool
}
