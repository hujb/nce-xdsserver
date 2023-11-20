package model

type Instance struct {
	InstanceId  string
	Port        uint32
	Ip          string
	Weight      uint32
	Metadata    map[string]string
	ClusterName string
	ServiceName string
	Enabled     bool
	// 对应返回报文的valid
	Healthy   bool
	Ephemeral bool
}

type Service struct {
	GroupName        string
	Name             string
	ProtectThreshold float32
	AppName          string
	Metadata         map[string]string
}
