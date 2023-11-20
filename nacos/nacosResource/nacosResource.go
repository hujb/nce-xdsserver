package nacosResource

// ServiceInstanceList 针对2023.11.16获取nacos数据的接口调整新方案新增，评估中
type ServiceInstanceList struct {
	Name                     string           `json:"name"`
	GroupName                string           `json:"groupName"`
	Clusters                 string           `json:"clusters"`
	CacheMillis              int              `json:"cacheMillis"`
	Hosts                    []*NacosInstance `json:"hosts"`
	LastRefTime              int64            `json:"lastRefTime"`
	Checksum                 string           `json:"checksum"`
	AllIPs                   bool             `json:"allIPs"`
	ReachProtectionThreshold bool             `json:"reachProtectionThreshold"`
	Valid                    bool             `json:"valid"`
}

type ServiceClusterInstanceDetail struct {
	Namespace        string                      `json:"namespace"`
	ServiceName      string                      `json:"serviceName"`
	GroupName        string                      `json:"groupName"`
	ClusterMap       map[string]ClusterMapDetail `json:"clusterMap"`
	Metadata         map[string]string           `json:"metadata"`
	ProtectThreshold float64                     `json:"protectThreshold"`
	Selector         interface{}                 `json:"selector"`
	Ephemeral        interface{}                 `json:"ephemeral"`
}

type ClusterMapDetail struct {
	ClusterName   string            `json:"clusterName"`
	HealthChecker HealthChecker     `json:"HealthChecker"`
	Metadata      map[string]string `json:"metadata"`
	Hosts         []*NacosInstance  `json:"hosts"`
}

type ResquestType struct {
}

type NacosNamespace struct {
	NacosNamespaceList map[string]string
}

type NacosGroup struct {
	NacosGroupList map[string]string
}

type NacosService struct {
	Name                string `json:"name"`
	GroupName           string `json:"groupName"`
	ClusterCount        int32  `json:"clusterCount"`
	IpCount             int32  `json:"ipCount"`
	HealthyInstanceCont int32  `json:"healthyInstanceCount"`
	TriggerFlag         string `json:"triggerFlag"`
}

type InstanceMetadata struct {
	List  []NacosInstance `json:"list"`
	Count int             `json:"count"`
}

type NacosInstance struct {
	InstanceId  string  `json:"instanceId"`
	Ip          string  `json:"ip"`
	Port        uint32  `json:"port"`
	Weight      float32 `json:"weight"`
	Valid       bool    `json:"valid"`
	Healthy     bool    `json:"healthy"`
	Enabled     bool    `json:"enabled"`
	Ephemeral   bool    `json:"ephemeral"`
	ClusterName string  `json:"clusterName"`
	ServiceName string  `json:"serviceName"`
	//ClusterSyncId string            `json:"clusterSyncId"`
	Metadata map[string]string `json:"metadata"`
	//LastBeat string            `json:"lastBeat"`
	//Marked string `json:"marked"`
	//App    string `json:"app"`
	//SyncChecksumMap           string            `json:"syncChecksumMap"`
	InstanceHeartBeatInterval int    `json:"InstanceHeartBeatInterval"`
	InstanceHeartBeatTimeOut  int    `json:"instanceHeartBeatTimeOut"`
	IpDeleteTimeout           int    `json:"ipDeleteTimeout"`
	InstanceIdGenerator       string `json:"InstanceIdGenerator"`
}

type ServiceClusterDetail struct {
	Service  ServiceDetail   `json:"service"`
	Clusters []ClusterDetail `json:"clusters"`
}

type ServiceDetail struct {
	Name             string            `json:"name"`
	GroupName        string            `json:"groupName"`
	ProtectThreshold float32           `json:"ProtectThreshold"`
	Selector         Selector          `json:"selector"`
	Metadata         map[string]string `json:"metadata"`
}

type Selector struct {
	Type        string `json:"type"`
	ContextType string `json:"ContextType"`
}

type ClusterDetail struct {
	ServiceName      string            `json:"ServiceName"`
	Name             string            `json:"name"`
	HealthChecker    HealthChecker     `json:"HealthChecker"`
	DefaultPort      int32             `json:"DefaultPort"`
	DefaultCheckPort int32             `json:"DefaultCheckPort"`
	UseIPPort4Check  bool              `json:"UseIPPort4Check"`
	Metadata         map[string]string `json:"metadata"`
}

type HealthChecker struct {
	Type string `json:"type"`
}
