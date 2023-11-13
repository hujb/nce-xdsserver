package util

import (
	"errors"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nce/nce-xdsserver/common/constant"
	"github.com/nce/nce-xdsserver/model"
	"github.com/nce/nce-xdsserver/nacos/nacosResource"
	mcp_v1alpha1 "istio.io/api/mcp/v1alpha1"
	"istio.io/api/networking/v1alpha3"
	"regexp"
	"strings"
	"time"
)

const VALID_DEFAULT_GROUP_NAME = "DEFAULT-GROUP"

func BuildServiceNameForServiceEntry(serviceDetail *nacosResource.ServiceClusterInstanceDetail, clusterName string, namespace string) string {
	var group string
	if constant.DEFAULT_GROUP != serviceDetail.GroupName {
		group = serviceDetail.GroupName
	} else {
		group = VALID_DEFAULT_GROUP_NAME
	}
	return serviceDetail.ServiceName + "." + clusterName + "." + group + "." + namespace
}

func BuildServiceEntry(svcName string, domainSuffix string, istioService *model.IstioService) *model.ServiceEntryWrapper {
	if istioService.Hosts == nil {
		return nil
	}
	var port uint32 = 0
	var protocol = "http"
	var hostname = svcName

	se := &v1alpha3.ServiceEntry{
		Hosts:      []string{hostname + "." + domainSuffix},
		Resolution: v1alpha3.ServiceEntry_STATIC,
		Location:   v1alpha3.ServiceEntry_MESH_INTERNAL,
	}

	for _, instance := range istioService.Hosts {
		if port == 0 {
			port = instance.Port
		}

		if (instance.Metadata)["protocol"] != "" {
			protocol = (instance.Metadata)["protocol"]
			if "triple" == protocol || "tri" == protocol {
				protocol = "grpc"
			}
		}
		metaHostname := (instance.Metadata)[constant.ISTIO_HOSTNAME]
		if metaHostname != "" {
			hostname = metaHostname
		}

		if !instance.Healthy || !instance.Enabled {
			continue
		}

		metadata := make(map[string]string)
		if instance.ClusterName != "" {
			metadata["cluster"] = instance.ClusterName
		}

		keyRegex, err := regexp.Compile(constant.VALID_LABEL_KEY_FORMAT)
		if err != nil {
			errors.New("无效的正则表达式:" + err.Error())
		}
		valueRegex, err := regexp.Compile(constant.VALID_LABEL_VALUE_FORMAT)
		if err != nil {
			errors.New("无效的正则表达式:" + err.Error())
		}

		for key, value := range instance.Metadata {
			if !keyRegex.MatchString(key) {
				continue
			}
			if !valueRegex.MatchString(value) {
				continue
			}
			metadata[key] = value
		}
		endpoint := &v1alpha3.WorkloadEntry{
			Labels:  metadata,
			Address: instance.Ip,
			Weight:  uint32(instance.Weight),
			Ports:   map[string]uint32{protocol: instance.Port},
		}
		se.Endpoints = append(se.Endpoints, endpoint)
	}

	servicePort := &v1alpha3.ServicePort{
		Number:   port,
		Protocol: strings.ToUpper(protocol),
		Name:     protocol,
	}
	se.Ports = append(se.Ports, servicePort)

	metadata := &mcp_v1alpha1.Metadata{
		Name:        istioService.Namespace + "/" + svcName,
		CreateTime:  &timestamp.Timestamp{Seconds: time.Now().Unix()},
		Labels:      map[string]string{"registryType": "nacos"},
		Version:     string(istioService.Revision),
		Annotations: map[string]string{"virtual": "1"},
	}

	return model.NewServiceEntryWrapper(se, metadata)
}
