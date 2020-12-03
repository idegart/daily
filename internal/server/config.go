package server

type Config struct {
	BindAddr string
}

func NewConfig(addr string) *Config {
	return &Config{
		BindAddr: ":" + addr,
	}
}