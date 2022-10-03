package webserver

const (
	DefaultHttpListenAddr = "127.0.0.1:8080"
)

type Config struct {
	DevMode        bool
	TLS            bool
	ConfigDir      string
	HttpListenAddr string
}

func NewConfig(configDir string) *Config {
	return &Config{
		DevMode:        true,
		TLS:            true,
		ConfigDir:      configDir,
		HttpListenAddr: DefaultHttpListenAddr,
	}
}

func (c *Config) Verify() error {

	return nil
}
