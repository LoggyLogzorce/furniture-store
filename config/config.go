package config

type Server struct {
	Port    string
	Timeout string
}

type Service struct {
	Method string
	Url    string
}

type AccessExceptions struct {
	List []string
}

type Config struct {
	Server
	Api              map[string]Service
	AccessExceptions `yaml:"accessExceptions"`
}

var config Config

func Get() *Config {
	return &config
}
