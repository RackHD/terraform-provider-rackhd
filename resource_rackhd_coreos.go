package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jfrey/go-rackhd"
	"github.com/jfrey/terraform-provider-rackhd/types"
)

func resourceRackHDCoreos() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackHDCoreosCreate,
		Read:   resourceRackHDCoreosRead,
		Update: resourceRackHDCoreosUpdate,
		Delete: resourceRackHDCoreosDelete,

		Schema: map[string]*schema.Schema{
			// Optional specification for a Node ID directly.
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			// Compute Host Configuration
			"hostname": &schema.Schema{
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
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"etcdToken": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRackHDCoreosCreate(d *schema.ResourceData, meta interface{}) error {
	resourceID, err := resourceIdentify(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*rackhd.Client)

	workflow, err := composeCoreosOptions(d)
	if err != nil {
		return err
	}

	err = client.NodeRunWorkflow(client.NodePostWorkflow, client.GetWorkflowByID, resourceID, workflow)
	if err != nil {
		return err
	}

	d.SetId(resourceID)

	return resourceRackHDCoreosRead(d, meta)
}

func resourceRackHDCoreosRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackHDCoreosUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRackHDCoreosRead(d, meta)
}

func resourceRackHDCoreosDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceCheckin(d, meta)
}

func composeCoreosOptions(d *schema.ResourceData) (interface{}, error) {
	workflow := types.CoreosWorkflow{
		Name: "Graph.InstallCoreOS",
		Options: types.CoreosOptions{
			InstallCoreos: types.InstallCoreos{
				Username:  d.Get("user").(string),
				Password:  d.Get("password").(string),
				Hostname:  d.Get("hostname").(string),
				SshKey:    d.Get("sshkey").(string),
				EtcdToken: d.Get("etcdToken").(string),
			},
		},
	}

	return workflow, nil
}
