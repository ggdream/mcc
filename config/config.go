package config

import "gopkg.in/yaml.v3"

var c *Config

type Config struct {
	DB   string `yaml:"db"`
	Logs string `yaml:"logs"`

	RunsBaseDir   string `yaml:"runsBaseDir"`
	ServerBaseDir string `yaml:"serverBaseDir"`
	StaticBaseDir string `yaml:"staticBaseDir"`

	Proxy Proxy `yaml:"proxy"`

	Gitea  Gitea  `yaml:"gitea"`
	Github Github `yaml:"github"`
	Gitlab Gitlab `yaml:"gitlab"`
}

type Proxy struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Gitea struct {
	Auth Auth `yaml:"auth"`
}

type Github struct {
	Auth Auth `yaml:"auth"`
}

type Gitlab struct {
	Auth Auth `yaml:"auth"`
}

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Init(data []byte) (err error) {
	err = yaml.Unmarshal(data, &c)
	return
}

func Get() *Config {
	return c
}
