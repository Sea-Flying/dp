package config

type DPConfig struct {
	Base      BaseConfig
	Cassandra CassandraConfig
	Consul    ConsulConfig
	Nomad     NomadConfig
	Gitlab    GitlabConfig
	AppPanel  AppPanelConfig
}

type BaseConfig struct {
	HttpPort string `yaml:"http-port" env:"BASE_HTTP_PORT" env-default:"7654"`
	HttpHost string `yaml:"http-host" env:"BASE_HTTP_HOST" env-default:""`
}

type CassandraConfig struct {
	HostsUrls string `yaml:"hosts-urls" env:"CASSANDRA_HOSTS_URLS" env-default:"127.0.0.1:9042"`
}

type ConsulConfig struct {
	ConsulApiUrl     string `yaml:"api-url" env:"CONSUL_ADDR" env-default:"http://10.0.0.152:8500"`
	ConsulDataCenter string `yaml:"datacenter" env:"CONSUL_DATACENTER" env-default:"vo-local"`
}

type NomadConfig struct {
	NomadApiUrl    string `yaml:"api-url" env:"NOMAD_ADDR" env-default:"http://10.0.0.152:4646"`
	NomadRegion    string `yaml:"region" env:"NOMAD_REGION" env-default:"vo-local"`
	NomadJobTplDir string `yaml:"job-tpl-dir" env:"NOMAD_JOB_TPL_DIR" env-default:"D:/tmp/jobtpldir"`
	NomadJobHclDir string `yaml:"job-hcl-dir" env:"NOMAD_JOB_HCL_DIR" env-default:"D:/tmp/jobhcldir"`
}

type GitlabConfig struct {
	Username string `yaml:"username" env:"GITLAB_USERNAME" env-default:"dp"`
	Token    string `yaml:"token" env:"GITLAB_TOKEN" env-default:"f4fc2yUEsVxdfa9zG9ev"`
}

type AppPanelConfig struct {
	WatcherAppsRegexpInclude string `yaml:"watcher-apps-regexp-include" env:"WATHCER_APPS_REGEXP_INCLUDE" env-default:"(?i)openvms.*|voerp.*"`
	WatcherAppsRegexpExclude string `yaml:"watcher-apps-regexp-exclude" env:"WATHCER_APPS_REGEXP_EXCLUDE" env-default:"(?i).*front.*"`
}
