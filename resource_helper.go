package main

import (
	"errors"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jfrey/go-rackhd"
)

var resourceCheckoutLock sync.Mutex

func resourceCheckout(d *schema.ResourceData, meta interface{}) (*rackhd.Node, error) {
	resourceCheckoutLock.Lock()
	defer resourceCheckoutLock.Unlock()

	client := meta.(*rackhd.Client)

	nodes, err := client.GetNodes()
	if err != nil {
		return nil, err
	}

	// Node selection looks for compute nodes which are not reserved.
	var node *rackhd.Node
	for _, current := range nodes {
		if current.Type == "compute" && !current.Terraform {
			node = &current
			break
		}
	}

	// If we couldn't find a Node we can't fulfill the resource request.
	if node == nil {
		return nil, errors.New("Unable to find eligible Node.")
	}

	// Patch the node to identify it's been reserved
	node.Terraform = true

	node, err = client.PatchNode(node.ID, node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func resourceCheckin(d *schema.ResourceData, meta interface{}) error {
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

	// TODO: Run a clean out workflow on the target Node.

	d.SetId("")

	return nil
}
