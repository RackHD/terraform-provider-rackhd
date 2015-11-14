package main

import (
	"errors"
	"sync"

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
	// Patch the node to identify it's been reserved
	client := meta.(*rackhd.Client)

	node, err := client.GetNode(d.Id())
	if err != nil {
		return err
	}

	// Patch the node to identify it's been un-reserved
	node.Terraform = false

	node, err = client.PatchNode(node.ID, node)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func updateResourceRackHDServerFromNode(d *schema.ResourceData, meta interface{}, node *rackhd.Node) error {
	client := meta.(*rackhd.Client)

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

	// Set the Node ID.
	d.SetId(node.ID)

	return nil
}

var selectResourceRackHDServerLock sync.Mutex

func selectResourceRackHDServer(d *schema.ResourceData, meta interface{}) (*rackhd.Node, error) {
	selectResourceRackHDServerLock.Lock()
	defer selectResourceRackHDServerLock.Unlock()

	client := meta.(*rackhd.Client)

	nodes, err := client.GetNodes()
	if err != nil {
		return nil, err
	}

	// Node selection looks for compute nodes which are not reserved.
	var node *rackhd.Node
	for _, current := range nodes {
		// TODO: add another field type for reservation.
		if current.Type == "compute" && !current.Terraform {
			node = &current
			break
		}
	}

	// If we couldn't find a Node we can't fulfill the resource request.
	if node == nil {
		return nil, errors.New("Unable to find eligible compute Node.")
	}

	// Patch the node to identify it's been reserved
	node.Terraform = true

	node, err = client.PatchNode(node.ID, node)
	if err != nil {
		return nil, err
	}

	return node, nil
}
