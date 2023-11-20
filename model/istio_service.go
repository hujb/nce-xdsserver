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
