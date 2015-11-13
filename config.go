package main

import (
	"log"
	"time"

	"github.com/jfrey/go-rackhd"
)

// Config contains the RackHD API configuration parameters
// required to connect the RackHD Go client to the API.
type Config struct {
	Host            string
	Port            int
	WorkflowTimeout time.Duration
}

// Client creates a RackHD API client which is
// utilized by the terraform.Provider.
func (c *Config) Client() (*rackhd.Client, error) {
	client := rackhd.Client{c.Host, c.Port, c.WorkflowTimeout}

	log.Printf("RackHD Client Configured: %s:%d\n", c.Host, c.Port)

	return &client, nil
}
