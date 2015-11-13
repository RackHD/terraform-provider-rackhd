package main

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jfrey/go-rackhd"
)

func resourceRackHDServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackHDServerCreate,
		Read:   resourceRackHDServerRead,
		Update: resourceRackHDServerUpdate,
		Delete: resourceRackHDServerDelete,

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

func resourceRackHDServerCreate(d *schema.ResourceData, meta interface{}) error {
	node, err := selectResourceRackHDServer(d, meta)
	if err != nil {
		return err
	}

	// TODO: Patch the selected node with the reservation field. This could be
	// part of the selectResourceRackHDServer function as well.

	// TODO: Execute Workflow

	// Update the resource properties.
	err = updateResourceRackHDServerFromNode(d, meta, node)
	if err != nil {
		return err
	}

	return resourceRackHDServerRead(d, meta)
}

func resourceRackHDServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*rackhd.Client)

	node, err := client.GetNode(d.Id())
	if err != nil {
		return err
	}

	return updateResourceRackHDServerFromNode(d, meta, node)
}

func resourceRackHDServerUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRackHDServerRead(d, meta)
}

func resourceRackHDServerDelete(d *schema.ResourceData, meta interface{}) error {
	// TODO: Run a clean out workflow on the target Node.

	d.SetId("")

	return nil
}

func updateResourceRackHDServerFromNode(d *schema.ResourceData, meta interface{}, node *rackhd.Node) error {
	client := meta.(*rackhd.Client)

	// Set computed properties from the Node document.
	d.Set("sku", node.Sku)
	d.Set("type", node.Type)

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

	// Set the Node ID.
	d.SetId(node.ID)

	return nil
}

func selectResourceRackHDServer(d *schema.ResourceData, meta interface{}) (*rackhd.Node, error) {
	client := meta.(*rackhd.Client)

	nodes, err := client.GetNodes()
	if err != nil {
		return nil, err
	}

	// Node selection looks for compute nodes which are not reserved.
	var selected *rackhd.Node
	for _, node := range nodes {
		// TODO: add another field type for reservation.
		if node.Type == "compute" {
			selected = &node
			break
		}
	}

	// If we couldn't find a Node we can't fulfill the resource request.
	if selected == nil {
		return nil, errors.New("Unable to find eligible compute Node.")
	}

	return selected, nil
}
