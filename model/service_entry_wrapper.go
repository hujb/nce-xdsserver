package model

import (
	mcp_v1alpha1 "istio.io/api/mcp/v1alpha1"
	"istio.io/api/networking/v1alpha3"
)

type ServiceEntryWrapper struct {
	se       *v1alpha3.ServiceEntry
	metadata *mcp_v1alpha1.Metadata
}

func NewServiceEntryWrapper(v1alpha3Se *v1alpha3.ServiceEntry, mcpMetadata *mcp_v1alpha1.Metadata) *ServiceEntryWrapper {
	return &ServiceEntryWrapper{se: v1alpha3Se, metadata: mcpMetadata}
}

func (w *ServiceEntryWrapper) GetServiceEntry() *v1alpha3.ServiceEntry {
	return w.se
}

func (w *ServiceEntryWrapper) GetMetadata() *mcp_v1alpha1.Metadata {
	return w.metadata
}
