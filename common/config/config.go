package config

type ServerConfig struct {
	Scheme      string //the nacos server scheme
	ContextPath string //the nacos server contextpath
	IpAddr      string //the nacos server address
	Port        uint64 //the nacos server port
}
