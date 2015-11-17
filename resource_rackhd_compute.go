package main

import (
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
			"workflow": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sku": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRackHDComputeCreate(d *schema.ResourceData, meta interface{}) error {
	node, err := resourceCheckout(d, meta)
	if err != nil {
		return err
	}

	d.SetId(node.ID)

	// TODO: Execute Workflow

	return resourceRackHDComputeRead(d, meta)
}

func resourceRackHDComputeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*rackhd.Client)

	node, err := client.GetNode(d.Id())
	if err != nil {
		return err
	}

	// Set computed properties from the Node document.
	d.Set("type", node.Type)

	// Get the friendly entry for the Sku if possible.
	if node.Sku != "" {
		sku, err := client.GetSku(node.Sku)
		if err == nil {
			d.Set("sku", sku.Name)
		} else {
			d.Set("sku", node.Sku)
		}
	}

	// Get the Lookup entries for the node for a computed IP Addresses property.
	lookups, err := client.QueryLookups(node.ID)
	if err != nil {
		return err
	}

	// Assign all valid IP Addresses
	var ipaddresses []string
	for _, lookup := range lookups {
		if lookup.IPAddress != "" {
			ipaddresses = append(ipaddresses, lookup.IPAddress)
		}
	}

	// Set the computed property for IP Addresses.
	d.Set("ipaddresses", ipaddresses)

	return nil
}

func resourceRackHDComputeUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRackHDComputeRead(d, meta)
}

func resourceRackHDComputeDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceCheckin(d, meta)
}
