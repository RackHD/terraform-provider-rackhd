package rackhd

import rhd "github.com/jfrey/go-rackhd"

// Config contains the RackHD API configuration parameters
// required to connect the RackHD Go client to the API.
type Config struct {
	Host string
	Port int
}

// Client creates a RackHD API client which is
// utilized by the terraform.Provider.
func (c *Config) Client() (*rhd.Client, error) {
	client := rhd.Client{c.Host, c.Port}

	return &client, nil
}
