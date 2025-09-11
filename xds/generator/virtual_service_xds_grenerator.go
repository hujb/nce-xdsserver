package generator

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nce/nce-xdsserver/common/resource"
	"github.com/nce/nce-xdsserver/log"
	"google.golang.org/protobuf/types/known/anypb"
	mcp_v1alpha1 "istio.io/api/mcp/v1alpha1"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"sync"
	"time"
)

type VirtualServiceXdsGenerator struct {
}

var singletonVirtualServiceXdsGenerator *VirtualServiceXdsGenerator

var onceVirtualService sync.Once

func GetVirtualServiceXdsGeneratorInstance() *VirtualServiceXdsGenerator {
	onceVirtualService.Do(func() {
		singletonVirtualServiceXdsGenerator = &VirtualServiceXdsGenerator{}
	})
	return singletonVirtualServiceXdsGenerator
}

func (s *VirtualServiceXdsGenerator) Generate(rs *resource.ResourceSnapshot) []*anypb.Any {
	httpRoutes := make([]*networkingv1alpha3.HTTPRoute, 0)
	httpRoute := &networkingv1alpha3.HTTPRoute{
		Match: []*networkingv1alpha3.HTTPMatchRequest{
			{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Prefix{
						Prefix: "/",
					},
				},
			},
		},
		Route: []*networkingv1alpha3.HTTPRouteDestination{
			{
				Destination: &networkingv1alpha3.Destination{
					Host: "reviews",
					Port: &networkingv1alpha3.PortSelector{
						Number: 9091,
					},
				},
			},
		},
	}
	httpRoutes = append(httpRoutes, httpRoute)
	vs := &networkingv1alpha3.VirtualService{
		Hosts:    []string{"my-virtual-service"},
		Http:     httpRoutes,
		ExportTo: []string{"*"},
	}
	result := make([]*anypb.Any, 0)
	toAny, err := ConvertVirtualServiceToAny(vs)
	if err != nil {
		log.Logger.Error(" convert virtualService to any failed, err: " + err.Error())
		return nil
	}
	result = append(result, toAny)
	return result
}

func ConvertVirtualServiceToAny(vs *networkingv1alpha3.VirtualService) (*anypb.Any, error) {
	// 将VirtualService对象编码为字节数组
	anyVs, err := anypb.New(vs)
	if err != nil {
		return nil, err
	}
	metadata := &mcp_v1alpha1.Metadata{
		Name:        "default/virtual-service",
		CreateTime:  &timestamp.Timestamp{Seconds: time.Now().Unix()},
		Labels:      map[string]string{"registryType": "apollo"},
		Version:     string(2),
		Annotations: map[string]string{"virtual": "2"},
	}
	//anyVs.TypeUrl = constant.VIRTUAL_SERVICE_PROTO
	r := &mcp_v1alpha1.Resource{Body: anyVs, Metadata: metadata}
	apb, _ := anypb.New(r)
	//apb.TypeUrl = constant.MCP_RESOURCE_PROTO

	//bytes, err := proto.Marshal(vs)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 创建一个anypb.Any对象
	//any, err := anypb.New(&anypb.Any{
	//	TypeUrl: constant.VIRTUAL_SERVICE_PROTO,
	//	Value:   bytes,
	//})
	////any := &anypb.Any{
	////	TypeUrl: constant.VIRTUAL_SERVICE_PROTO,
	////	Value:   bytes,
	////}
	//if err != nil {
	//	return nil, err
	//}

	return apb, nil
}
