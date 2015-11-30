package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

// VSphere struct holds the config for communicating with a VSphere Server.
type VSphere struct {
	Host     string
	User     string
	Password string
}

// Client creates a govmomi.Client for use in communicating with a VSphere Server.
func (c *VSphere) Client() (*govmomi.Client, error) {
	u, err := url.Parse("https://" + c.Host + "/sdk")
	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword(c.User, c.Password)

	ctx := context.TODO()

	timeout := time.NewTimer(time.Minute * 30).C
	retry := time.NewTicker(time.Second * 10).C

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("Timed Out Creating Client.")
		case <-retry:
			client, err := govmomi.NewClient(ctx, u, true)
			if err == nil {
				return client, nil
			}
		}
	}
}

// AddClusterToDefaultDatacenter adds the given cluster to the default datacenter.
func AddClusterToDatacenter(c *govmomi.Client, cluster string, datacenter string) error {
	ctx := context.TODO()
	finder := find.NewFinder(c.Client, false)

	dc, err := finder.Datacenter(ctx, datacenter)
	if err != nil {
		return err
	}

	folders, err := dc.Folders(ctx)
	if err != nil {
		return err
	}

	parent := folders.HostFolder

	_, err = parent.CreateCluster(ctx, cluster, types.ClusterConfigSpecEx{})
	if err != nil {
		return err
	}

	return nil
}

// AddDatacenter adds a datacenter to VCenter.
func AddDatacenter(c *govmomi.Client, name string) error {
	ctx := context.TODO()
	rootFolder := object.NewRootFolder(c.Client)

	_, err := rootFolder.CreateDatacenter(ctx, name)
	if err != nil {
		return err
	}

	return nil
}

func AddHostToCluster(c *govmomi.Client, datacenter string, cluster string, host string, user string, password string) error {
	ctx := context.Background()

	finder := find.NewFinder(c.Client, true)

	dc, err := finder.Datacenter(ctx, datacenter)
	if err != nil {
		return err
	}

	finder.SetDatacenter(dc)

	ccr, err := finder.ClusterComputeResource(ctx, cluster)
	if err != nil {
		return err
	}

	spec := types.HostConnectSpec{}
	spec.HostName = host
	spec.UserName = user
	spec.Password = password

	task, err := ccr.AddHost(ctx, spec, true, nil, nil)
	if err != nil {
		return err
	}

	_, err = task.WaitForResult(ctx, nil)
	if err != nil {
		if f, ok := err.(types.HasFault); ok {
			switch fault := f.Fault().(type) {
			case *types.SSLVerifyFault:
				spec.SslThumbprint = fault.Thumbprint

				task, err = ccr.AddHost(ctx, spec, true, nil, nil)
				if err != nil {
					return err
				}

				_, err = task.WaitForResult(ctx, nil)
				if err != nil {
					return err
				}
			}

		}

		return err
	}

	return nil
}
