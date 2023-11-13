package constant

const (
	SERVICE_ENTRY_PROTO_PACKAGE = "networking.istio.io/v1alpha3/ServiceEntry"
	//SERVICE_ENTRY_TYPE          = "networking.istio.io/v1alpha3/ServiceEntry"

	SERVICE_ENTRY_TYPE = "istio.networking.v1alpha3.ServiceEntry"

	MESH_CONFIG_TYPE = "core/v1alpha1/MeshConfig"

	MCP_PREFIX = "istio/"

	DEFAULT_GROUP = "DEFAULT_GROUP"

	NACOS_ISTIO_DOMAIN_SUFFIX = "nacos"

	ISTIO_HOSTNAME = "istio.hostname"

	VALID_LABEL_KEY_FORMAT = "^([a-zA-Z0-9](?:[-a-zA-Z0-9]*[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[-a-zA-Z0-9]*[a-zA-Z0-9])?)*/)?((?:[A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$"

	VALID_LABEL_VALUE_FORMAT = "^((?:[A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"

	API_TYPE_PREFIX = "type.googleapis.com/"

	SERVICE_ENTRY_PROTO = API_TYPE_PREFIX + SERVICE_ENTRY_TYPE

	MCP_RESOURCE_PROTO = API_TYPE_PREFIX + "istio.mcp.v1alpha1.Resource"
)
