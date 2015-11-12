package main

import (
	"log"

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
			"os": &schema.Schema{
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
		},
	}
}

/*
{
	AutoDiscover:false
	Catalogs:[]
	CreatedAt:2015-11-11T18:24:13.166Z
	ID:564387cd64fd10000546194d
	Identifiers:[52:54:be:ef:87:03]
	Name:52:54:be:ef:87:03
	ObmSettings:[]
	Relations:[]
	Sku:
	Type:compute
	UpdatedAt:2015-11-11T18:36:08.635Z
	Workflows:[]
}
*/

func resourceRackHDServerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("RackHD Create Resource")

	client := meta.(*rackhd.Client)

	nodes, err := client.GetNodes()
	if err != nil {
		log.Printf("Error getting Nodes: %s\n", err)
		return err
	}

	var selected *rackhd.Node

	for _, node := range nodes {
		if node.Type == "compute" {
			selected = &node
			break
		}
	}

	if selected == nil {
		log.Println("Unable to find eligible compute Node.")
	}

	log.Printf("Selected Node: %+v\n", selected)
	log.Printf("%s", selected.ID)

	d.Set("sku", selected.Sku)
	d.Set("type", selected.Type)

	d.SetId(selected.ID)

	//
	// app := d.Get("app").(string)
	// opts := heroku.AddonCreateOpts{Plan: d.Get("plan").(string)}
	//
	// if v := d.Get("config"); v != nil {
	// 	config := make(map[string]string)
	// 	for _, v := range v.([]interface{}) {
	// 		for k, v := range v.(map[string]interface{}) {
	// 			config[k] = v.(string)
	// 		}
	// 	}
	//
	// 	opts.Config = &config
	// }
	//
	// log.Printf("[DEBUG] Addon create configuration: %#v, %#v", app, opts)
	// a, err := client.AddonCreate(app, opts)
	// if err != nil {
	// 	return err
	// }
	//
	// d.SetId(a.ID)
	// log.Printf("[INFO] Addon ID: %s", d.Id())

	return resourceRackHDServerRead(d, meta)
}

func resourceRackHDServerRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("RackHD Read Resource")

	client := meta.(*rackhd.Client)

	node, err := client.GetNode(d.Id())
	if err != nil {
		log.Printf("Error getting Node: %s\n", err)
		return err
	}

	log.Printf("Read Node: %+v\n", node)

	d.Set("sku", node.Sku)
	d.Set("type", node.Type)

	// client := meta.(*heroku.Service)
	//
	// addon, err := resourceRackHDServerRetrieve(
	// 	d.Get("app").(string), d.Id(), client)
	// if err != nil {
	// 	return err
	// }
	//
	// // Determine the plan. If we were configured without a specific plan,
	// // then just avoid the plan altogether (accepting anything that
	// // Heroku sends down).
	// plan := addon.Plan.Name
	// if v := d.Get("plan").(string); v != "" {
	// 	if idx := strings.IndexRune(v, ':'); idx == -1 {
	// 		idx = strings.IndexRune(plan, ':')
	// 		if idx > -1 {
	// 			plan = plan[:idx]
	// 		}
	// 	}
	// }
	//
	// d.Set("name", addon.Name)
	// d.Set("plan", plan)
	// d.Set("provider_id", addon.ProviderID)
	// if err := d.Set("config_vars", addon.ConfigVars); err != nil {
	// 	return err
	// }

	return nil
}

func resourceRackHDServerUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("RackHD Update Resource")

	return resourceRackHDServerRead(d, meta)
}

func resourceRackHDServerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("RackHD Delete Resource")

	client := meta.(*rackhd.Client)

	node, err := client.GetNode(d.Id())
	if err != nil {
		log.Printf("Error getting Node: %s\n", err)
		return err
	}

	log.Printf("Read Node: %+v\n", node)

	d.SetId("")
	return nil
}
