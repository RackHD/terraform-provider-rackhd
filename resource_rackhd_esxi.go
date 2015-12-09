package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jfrey/go-rackhd"
)

func resourceRackHDEsxi() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackHDEsxiCreate,
		Read:   resourceRackHDEsxiRead,
		Update: resourceRackHDEsxiUpdate,
		Delete: resourceRackHDEsxiDelete,

		Schema: map[string]*schema.Schema{
			// Optional specification for a Node ID directly.
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			// VCenter Configuration
			"vcenter_host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GOVC_URL", nil),
			},
			"vcenter_user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GOVC_USERNAME", nil),
			},
			"vcenter_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_cluster": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			// ESXI Host Configuration
			"esxi_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_interface": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_dns_server": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_root_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"esxi_uid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceRackHDEsxiCreate(d *schema.ResourceData, meta interface{}) error {
	var resourceID string

	// If node is set we'll try to use that as the resource ID, otherwise we'll
	// reserve a resource using resourceCheckout.
	if v, ok := d.GetOk("node"); ok {
		resourceID = v.(string)
	} else {
		node, err := resourceCheckout(d, meta)
		if err != nil {
			return err
		}

		resourceID = node.ID
	}

	client := meta.(*rackhd.Client)

	// TODO: Run esxi workflow
	user := rackhd.User{
		d.Get("esxi_user").(string),
		d.Get("esxi_password").(string),
		d.Get("esxi_uid").(int),
	}
	users := rackhd.Users{user}
	net := rackhd.NetworkDevice{
		d.Get("esxi_interface").(string),
		rackhd.Ipv4{
			d.Get("esxi_gateway").(string),
			d.Get("esxi_ip").(string),
			d.Get("esxi_subnet").(string),
		},
	}
	nets := rackhd.NetworkDevices{net}
	dns := []string{
		d.Get("esxi_dns_server").(string),
	}
	domain := d.Get("esxi_domain").(string)
	hostname := d.Get("esxi_hostname").(string)
	rootpass := d.Get("esxi_root_password").(string)
	version := d.Get("esxi_version").(string)
	install := rackhd.InstallOsStruct{
		dns,
		domain,
		hostname,
		rootpass,
		version,
		nets,
		users,
		d.Get("vcenter_host").(string),
		d.Get("esxi_subnet").(string),
		d.Get("esxi_gateway").(string),
	}
	defaults := rackhd.WorkflowOptionDefaults{
		InstallOs: install,
	}
	body := rackhd.WorkflowRequest{
		Name:    "Graph.InstallEsxvCenter",
		Options: defaults,
	}

	err := client.NodeRunWorkflow(client.NodePostWorkflow, client.GetWorkflowByID, resourceID, body)
	if err != nil {
		return err
	}

	vsphere := VSphere{}

	if v, ok := d.GetOk("vcenter_host"); ok {
		vsphere.Host = v.(string)
	}

	if v, ok := d.GetOk("vcenter_user"); ok {
		vsphere.User = v.(string)
	}

	if v, ok := d.GetOk("vcenter_password"); ok {
		vsphere.Password = v.(string)
	}

	vmomi, err := vsphere.Client()
	if err != nil {
		return err
	}

	err = AddDatacenter(
		vmomi,
		d.Get("vcenter_datacenter").(string),
	)
	if err != nil {
		return err
	}

	err = AddClusterToDatacenter(
		vmomi,
		d.Get("vcenter_cluster").(string),
		d.Get("vcenter_datacenter").(string),
	)
	if err != nil {
		return err
	}

	err = AddHostToCluster(
		vmomi,
		d.Get("vcenter_datacenter").(string),
		d.Get("vcenter_cluster").(string),
		d.Get("esxi_ip").(string),
		"root",
		d.Get("esxi_root_password").(string),
	)
	if err != nil {
		return err
	}

	d.SetId(resourceID)

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
