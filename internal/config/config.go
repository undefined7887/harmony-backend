package config

type Centrifugo struct {
	ApiKey     string `yaml:"api_key"`
	ApiAddress string `yaml:"api_address"`
}

type Mongo struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
