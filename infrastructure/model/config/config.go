package config

type DPConfig struct {
	Base      BaseConfig
	Cassandra CassandraConfig
	Nomad     NomadConfig
	Gitlab    GitlabConfig
}

type BaseConfig struct {
	HttpPort string `yaml:"http-port" env:"BASE_HTTP_PORT" `
	HttpHost string `yaml:"http-host" env:"BASE_HTTP_HOST" `
}

type CassandraConfig struct {
	HostsUrls string `yaml:"hosts-urls" env:"CASSANDRA_HOSTS_URLS" env-description:"connection urls of Cassandra instances, comma-separated"`
}

type NomadConfig struct {
	NomadApiUrl    string `yaml:"api-url" env:"NOMAD_ADDR"`
	NomadRegion    string `yaml:"region" env:"NOMAD_REGION"`
	NomadJobTplDir string `yaml:"job-tpl-dir" env:"NOMAD_JOB_TPL_DIR"`
	NomadJobHclDir string `yaml:"job-hcl-dir" env:"NOMAD_JOB_HCL_DIR"`
}

type GitlabConfig struct {
	Username string `yaml:"username" env:"GITLAB_USERNAME"`
	Token    string `yaml:"token" env:"GITLAB_TOKEN"`
}
