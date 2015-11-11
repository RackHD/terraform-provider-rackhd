package rackhd

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("RACKHD_HOST", nil),
			},

			"port": &schema.Schema{
				Type:        schema.TypeInteger,
				Optional:    true,
				DefaultFunc: trueschema.EnvDefaultFunc("RACKHD_PORT", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rackhd_server": resourceRackHDServer()
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Email:  d.Get("email").(string),
		APIKey: d.Get("api_key").(string),
	}

	log.Println("[INFO] Initializing Heroku client")
	return config.Client()
}
