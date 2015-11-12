package main

import (
	"strconv"

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
		},

		ResourcesMap: map[string]*schema.Resource{
			"rackhd_server": resourceRackHDServer(),
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

	config := Config{
		Host: d.Get("host").(string),
		Port: port,
	}

	return config.Client()
}
