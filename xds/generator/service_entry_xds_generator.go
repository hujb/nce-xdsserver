package generator

import (
	"github.com/nce/nce-xdsserver/common/constant"
	"github.com/nce/nce-xdsserver/common/resource"
	"google.golang.org/protobuf/types/known/anypb"
	mcp_v1alpha1 "istio.io/api/mcp/v1alpha1"
	"sync"
)

type ServiceEntryXdsGenerator struct {
}

var singletonServiceEntryXdsGenerator *ServiceEntryXdsGenerator

var once sync.Once

func GetInstance() *ServiceEntryXdsGenerator {
	once.Do(func() {
		singletonServiceEntryXdsGenerator = &ServiceEntryXdsGenerator{}
	})
	return singletonServiceEntryXdsGenerator
}

func (s *ServiceEntryXdsGenerator) Generate(rs *resource.ResourceSnapshot) []*anypb.Any {
	resources := make([]*mcp_v1alpha1.Resource, 0)
	for _, serviceEntryWrapper := range rs.GetServiceEntries() {
		metadata := serviceEntryWrapper.GetMetadata()
		serviceEntry := serviceEntryWrapper.GetServiceEntry()
		any, _ := anypb.New(serviceEntry)
		any.TypeUrl = constant.SERVICE_ENTRY_PROTO
		resource := &mcp_v1alpha1.Resource{Body: any, Metadata: metadata}
		resources = append(resources, resource)
	}
	result := make([]*anypb.Any, 0)
	for _, r := range resources {
		apb, _ := anypb.New(r)
		apb.TypeUrl = constant.MCP_RESOURCE_PROTO
		result = append(result, apb)
	}
	return result
}
