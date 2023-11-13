package model

type ServiceEntry struct {
	Hosts                []string                 `protobuf:"bytes,1,rep,name=hosts,proto3" json:"hosts,omitempty"`
	Addresses            []string                 `protobuf:"bytes,2,rep,name=addresses,proto3" json:"addresses,omitempty"`
	Ports                []*Port                  `protobuf:"bytes,3,rep,name=ports,proto3" json:"ports,omitempty"`
	Location             ServiceEntry_Location    `protobuf:"varint,4,opt,name=location,proto3,enum=istio.networking.v1alpha3.ServiceEntry_Location" json:"location,omitempty"`
	Resolution           ServiceEntry_Resolution  `protobuf:"varint,5,opt,name=resolution,proto3,enum=istio.networking.v1alpha3.ServiceEntry_Resolution" json:"resolution,omitempty"`
	Endpoints            []*ServiceEntry_Endpoint `protobuf:"bytes,6,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
	ExportTo             []string                 `protobuf:"bytes,7,rep,name=export_to,json=exportTo,proto3" json:"export_to,omitempty"`
	SubjectAltNames      []string                 `protobuf:"bytes,8,rep,name=subject_alt_names,json=subjectAltNames,proto3" json:"subject_alt_names,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

type Port struct {
	Number   uint32 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Protocol string `protobuf:"bytes,2,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Name     string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	//XXX_NoUnkeyedLiteral struct{} `json:"-"`
	//XXX_unrecognized     []byte   `json:"-"`
	//XXX_sizecache        int32    `json:"-"`
}

type ServiceEntry_Location int32

const (
	ServiceEntry_MESH_EXTERNAL ServiceEntry_Location = 0
	ServiceEntry_MESH_INTERNAL ServiceEntry_Location = 1
)

type ServiceEntry_Resolution int32

const (
	ServiceEntry_NONE   ServiceEntry_Resolution = 0
	ServiceEntry_STATIC ServiceEntry_Resolution = 1
	ServiceEntry_DNS    ServiceEntry_Resolution = 2
)

type ServiceEntry_Endpoint struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Ports   uint32 `protobuf:"bytes,2,rep,name=ports,proto3" json:"ports,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// 使用labels存储dmf元数据
	Labels               map[string]string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Network              string            `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
	Locality             string            `protobuf:"bytes,5,opt,name=locality,proto3" json:"locality,omitempty"`
	Weight               uint32            `protobuf:"varint,6,opt,name=weight,proto3" json:"weight,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}
