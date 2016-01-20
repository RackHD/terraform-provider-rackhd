package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jfrey/go-rackhd"
	"github.com/skunkworxs/terraform-provider-rackhd/types"
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
				Optional: true,
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

	workflow, err := composeWorkflow(d)
	if err != nil {
		return err
	}

	err = client.NodeRunWorkflow(client.NodePostWorkflow, client.GetWorkflowByID, resourceID, workflow)
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

func composeWorkflow(d *schema.ResourceData) (interface{}, error) {
	if v, ok := d.GetOk("os"); ok {
		// We're converting a known moniker for the compute OS into the RackHD
		// data required to execute the proper workflow.  This should be cleaned
		// up by unifying the install parameters across OS's in RackHD.
		switch v.(string) {
		case "centos-6.5":
			return composeCentosOptions(d, "6.5")
		case "centos-7.0":
			return composeCentosOptions(d, "7")
		case "ubuntu-trusty":
			return composeUbuntuOptions(d, "install-trusty.pxe")
		case "ubuntu-utopic":
			return composeUbuntuOptions(d, "install-utopic.pxe")
		}

		return nil, fmt.Errorf("Unsupported OS Requested: %s\n", v.(string))
	}

	return nil, fmt.Errorf("Required paramter 'os' not specified.")
}

func composeCentosOptions(d *schema.ResourceData, version string) (interface{}, error) {
	workflow := types.CentosWorkflow{
		Name: "Graph.InstallCentOS",
		Options: types.CentosOptions{
			InstallOS: types.InstallOS{
				Domain:       d.Get("domain").(string),
				Hostname:     d.Get("hostname").(string),
				RootPassword: d.Get("root_password").(string),
				Users: []types.User{
					types.User{
						Name:     d.Get("user").(string),
						Password: d.Get("password").(string),
						UID:      d.Get("uid").(int),
					},
				},
				Version: version,
			},
		},
	}

	return workflow, nil
}

func composeUbuntuOptions(d *schema.ResourceData, profile string) (interface{}, error) {
	workflow := types.UbuntuWorkflow{
		Name: "Graph.InstallUbuntu",
		Options: types.UbuntuOptions{
			InstallUbuntu: types.InstallUbuntu{
				Comport:       "ttyS0",
				CompletionUri: "renasar-ansible.pub",
				Domain:        d.Get("domain").(string),
				Hostname:      d.Get("hostname").(string),
				Password:      d.Get("password").(string),
				Profile:       profile,
				UID:           d.Get("uid").(int),
				Username:      d.Get("user").(string),
			},
		},
	}

	return workflow, nil
}
