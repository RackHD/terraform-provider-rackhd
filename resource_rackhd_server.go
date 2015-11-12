package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
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
		},
	}
}

func resourceRackHDServerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("RackHD Create Resource")

	// client := meta.(*rackhd.Client)
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
	// client := meta.(*heroku.Service)
	//
	// app := d.Get("app").(string)
	//
	// if d.HasChange("plan") {
	// 	ad, err := client.AddonUpdate(
	// 		app, d.Id(), heroku.AddonUpdateOpts{Plan: d.Get("plan").(string)})
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	// Store the new ID
	// 	d.SetId(ad.ID)
	// }

	return resourceRackHDServerRead(d, meta)
}

func resourceRackHDServerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("RackHD Delete Resource")
	// client := meta.(*heroku.Service)
	//
	// log.Printf("[INFO] Deleting Addon: %s", d.Id())
	//
	// // Destroy the app
	// err := client.AddonDelete(d.Get("app").(string), d.Id())
	// if err != nil {
	// 	return fmt.Errorf("Error deleting addon: %s", err)
	// }

	d.SetId("")
	return nil
}
