package rackhd

// Config contains the RackHD API configuration parameters
// required to connect the RackHD Go client to the API.
type Config struct {
	Host string
	Port integer
}

// Client creates a RackHD API client which is
// utilized by the terraform.Provider.
func (c *Config) Client() (*heroku.Service, error) {
	return c, nil
}
