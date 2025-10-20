package config

type ClusterConfig struct {
	Type    string
	Servers []string
}

type Config struct {
	Port          int32
	ClusterPort   int32
	AdvertisePort int32
	BindAddress   string
	Cluster       ClusterConfig
}
