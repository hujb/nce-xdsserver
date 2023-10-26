package util

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"gitlab2.psbc.com/ecn/ecn-xdsserver/common/constant"
	"gitlab2.psbc.com/ecn/ecn-xdsserver/model"
	mcp_v1alpha1 "istio.io/api/mcp/v1alpha1"
	"istio.io/api/networking/v1alpha3"
	"time"
)

const VALID_DEFAULT_GROUP_NAME = "DEFAULT-GROUP"

func BuildServiceNameForServiceEntry(service *model.Service) string {
	var group string
	if constant.DEFAULT_GROUP != service.GetGroup() {
		group = service.GetGroup()
	} else {
		group = VALID_DEFAULT_GROUP_NAME
	}
	return service.GetName() + "." + group + "." + service.GetNamespace()
}

func BuildServiceEntry(svcName string) *model.ServiceEntryWrapper {
	port := &v1alpha3.ServicePort{
		Number:   8080,
		Protocol: "HTTP",
		Name:     "http",
	}
	metadata := &mcp_v1alpha1.Metadata{
		Name:       "nacos/test",
		CreateTime: &timestamp.Timestamp{Seconds: time.Now().Unix()},
		Labels:     map[string]string{"hello": "test", "pa": "true"},
		Version:    "1111415485643",
	}
	se := &v1alpha3.ServiceEntry{
		Hosts:      []string{svcName + ".nacos"},
		Resolution: v1alpha3.ServiceEntry_STATIC,
		Location:   v1alpha3.ServiceEntry_MESH_INTERNAL,
		Ports:      []*v1alpha3.ServicePort{port},
	}
	return model.NewServiceEntryWrapper(se, metadata)
}
