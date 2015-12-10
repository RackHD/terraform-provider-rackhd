package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jfrey/go-rackhd"
)

func resourceRackHDCompute() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackHDComputeCreate,
		Read:   resourceRackHDComputeRead,
		Update: resourceRackHDComputeUpdate,
		Delete: resourceRackHDComputeDelete,

		Schema: map[string]*schema.Schema{
			// Optional specification for a Node ID directly.
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			// Compute Host Configuration
			"os": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"interface": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_server": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"root_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"uid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceRackHDComputeCreate(d *schema.ResourceData, meta interface{}) error {
	resourceID, err := resourceIdentify(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*rackhd.Client)

	body, err := composeWorkflowOptions(d)

	err = client.NodeRunWorkflow(client.NodePostWorkflow, client.GetWorkflowByID, resourceID, *body)
	if err != nil {
		return err
	}

	d.SetId(resourceID)

	return resourceRackHDComputeRead(d, meta)
}

func resourceRackHDComputeRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackHDComputeUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRackHDComputeRead(d, meta)
}

func resourceRackHDComputeDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceCheckin(d, meta)
}

func composeWorkflowOptions(d *schema.ResourceData) (*rackhd.WorkflowRequest, error) {
	if v, ok := d.GetOk("os"); ok {
		switch v.(string) {
		case "centos-6.5":
			return composeCentos65Options(d)
		}

		return nil, fmt.Errorf("Unsupported OS Requested: %s\n", v.(string))
	}

	return nil, fmt.Errorf("Required paramter 'os' not specified.")
}

func composeCentos65Options(d *schema.ResourceData) (*rackhd.WorkflowRequest, error) {
	user := rackhd.User{
		d.Get("user").(string),
		d.Get("password").(string),
		d.Get("uid").(int),
	}
	users := rackhd.Users{user}
	net := rackhd.NetworkDevice{
		d.Get("interface").(string),
		rackhd.Ipv4{
			d.Get("gateway").(string),
			d.Get("ip").(string),
			d.Get("subnet").(string),
		},
	}
	nets := rackhd.NetworkDevices{net}
	dns := []string{
		d.Get("dns_server").(string),
	}
	domain := d.Get("domain").(string)
	hostname := d.Get("hostname").(string)
	rootpass := d.Get("root_password").(string)
	version := "6.5"
	install := rackhd.InstallOsStruct{
		dns,
		domain,
		hostname,
		rootpass,
		version,
		nets,
		users,
		"",
		"",
		"",
	}
	defaults := rackhd.WorkflowOptionDefaults{
		InstallOs: install,
	}
	body := rackhd.WorkflowRequest{
		Name:    "Graph.InstallCentOS",
		Options: defaults,
	}

	return &body, nil
}
