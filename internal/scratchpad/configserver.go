package scratchpad

import "net/http"

type Configuration struct {
	PassPhrase string `yaml:"pass_phrase"`
}

type ConfigServer struct {
	Configuration *Configuration
}

func NewConfigServer(c *Configuration) *ConfigServer {
	return &ConfigServer{c}
}

func (c *ConfigServer) Start() error {
	mux := http.NewServeMux()
	mux.Handle("GET /api/register", handleClientRegistration(c.Configuration))

	return nil
}

func main() {
	c := &Configuration{PassPhrase: "This is a passphrase used to protect yourself"}
	s := NewConfigServer(c)
	_ = s.Start()
}
