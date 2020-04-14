package model

type DPConfig struct {
	Base      BaseConfig
	Cassandra CassandraConfig
	Nomad     NomadConfig
}

type BaseConfig struct {
	HttpPort string `yaml:"http-port" env:"BASE_HTTP_PORT" `
	HttpHost string `yaml:"http-host" env:"BASE_HTTP_HOST" `
}

type CassandraConfig struct {
	HostsUrls string `yaml:"hosts-urls" env:"CASSANDRA_HOSTS_URLS" env-description:"connection urls of Cassandra instances, comma-separated"`
}

type NomadConfig struct {
	NomadApiUrl string `yaml:"api-url" env:"NOMAD_API_URL"`
}
