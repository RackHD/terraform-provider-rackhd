package rackhd

import (
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
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKHD_PORT", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rackhd_server": resourceRackHDServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host: d.Get("host").(string),
		Port: d.Get("port").(int),
	}

	return config.Client()
}
