package main

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKHD_HOST", nil),
			},

			"port": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKHD_PORT", "8080"),
			},
			"workflow_timeout": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKHD_WORKFLOW_TIMEOUT", "7200"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rackhd_compute": resourceRackHDCompute(),
			"rackhd_esxi":    resourceRackHDEsxi(),
			"rackhd_coreos":  resourceRackHDCoreos(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// int doesn't appear to be available from terraform plan interrogation
	// so pulling it in as a string and converting it is a good short term option.
	port, err := strconv.Atoi(d.Get("port").(string))
	if err != nil {
		return nil, err
	}

	timeout, err := strconv.Atoi(d.Get("workflow_timeout").(string))
	if err != nil {
		return nil, err
	}

	config := Config{
		Host:            d.Get("host").(string),
		Port:            port,
		WorkflowTimeout: time.Second * time.Duration(timeout),
	}

	return config.Client()
}
