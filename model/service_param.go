package model

type GetAllServiceInfoParam struct {
	NameSpace     string
	GroupName     string
	PageNo        uint32
	PageSize      uint32
	WithInstances bool
}
