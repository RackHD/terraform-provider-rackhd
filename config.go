package main

import (
	"log"

	"github.com/jfrey/go-rackhd"
)

// Config contains the RackHD API configuration parameters
// required to connect the RackHD Go client to the API.
type Config struct {
	Host string
	Port int
}

// Client creates a RackHD API client which is
// utilized by the terraform.Provider.
func (c *Config) Client() (*rackhd.Client, error) {
	client := rackhd.Client{c.Host, c.Port}

	log.Printf("RackHD Client Configured: %s:%d\n", c.Host, c.Port)

	return &client, nil
}
