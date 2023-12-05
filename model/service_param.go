package model

import "fmt"

type QueryAllServiceInfoParam struct {
	NamespaceId   string
	GroupName     string
	WithInstances bool
	PageNo        uint32
	PageSize      uint32
}

func (param QueryAllServiceInfoParam) String() string {
	return fmt.Sprintf("NamespaceId: %s, GroupName: %s, WithInstances: %t, PageNo: %d, PageSize: %d",
		param.NamespaceId, param.GroupName, param.WithInstances, param.PageNo, param.PageSize)
}

type QueryAllInstanceInfoByServiceParam struct {
	ServiceName string
	GroupName   string
	NamespaceId string
	Clusters    string
	HealthyOnly bool
}

func (param QueryAllInstanceInfoByServiceParam) String() string {
	return fmt.Sprintf("ServiceName: %s, GroupName: %s, NamespaceId: %s, Clusters: %s, HealthyOnly: %t",
		param.NamespaceId, param.GroupName, param.NamespaceId, param.Clusters, param.HealthyOnly)
}
