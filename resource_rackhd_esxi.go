package main

import "github.com/hashicorp/terraform/helper/schema"

func resourceRackHDEsxi() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackHDEsxiCreate,
		Read:   resourceRackHDEsxiRead,
		Update: resourceRackHDEsxiUpdate,
		Delete: resourceRackHDEsxiDelete,

		Schema: map[string]*schema.Schema{
			"vcenter_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_pass": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRackHDEsxiCreate(d *schema.ResourceData, meta interface{}) error {
	node, err := resourceCheckout(d, meta)
	if err != nil {
		return err
	}

	// TODO: Run esxi workflow

	d.SetId(node.ID)

	return resourceRackHDEsxiRead(d, meta)
}

func resourceRackHDEsxiRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackHDEsxiUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRackHDEsxiRead(d, meta)
}

func resourceRackHDEsxiDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceCheckin(d, meta)
}
