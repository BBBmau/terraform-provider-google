package compute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeInstancesRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Elem:     &schema.Resource{Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstance().Schema)},
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeInstancesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	zone := d.Get("zone").(string)

	instances := make([]map[string]interface{}, 0)
	err = ListInstances(context.Background(), d, config, func(item interface{}) error {
		// Convert raw API response to compute.Instance
		itemMap := item.(map[string]interface{})
		itemJSON, err := json.Marshal(itemMap)
		if err != nil {
			return fmt.Errorf("Error marshaling instance: %s", err)
		}

		var instance compute.Instance
		if err := json.Unmarshal(itemJSON, &instance); err != nil {
			return fmt.Errorf("Error unmarshaling instance: %s", err)
		}

		// Create a temporary ResourceData for this instance
		instanceResource := ResourceComputeInstance()
		// Create an empty state to initialize the ResourceData
		emptyState := &terraform.InstanceState{
			Attributes: make(map[string]string),
		}
		tempData := instanceResource.Data(emptyState)
		tempData.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, instance.Name))
		tempData.Set("project", project)
		tempData.Set("zone", zone)
		tempData.Set("name", instance.Name)

		// Flatten the instance into the temporary ResourceData
		if err := flattenComputeInstance(tempData, config); err != nil {
			return fmt.Errorf("Error flattening instance: %s", err)
		}

		// Extract all values from the temporary ResourceData into a map
		// Use the data source schema to determine which fields to extract
		dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(instanceResource.Schema)
		instanceMap := make(map[string]interface{})
		for key := range dsSchema {
			if val, ok := tempData.GetOk(key); ok {
				// Only include non-nil values
				if val != nil {
					instanceMap[key] = val
				}
			}
		}

		instances = append(instances, instanceMap)
		return nil
	})
	if err != nil {
		return err
	}

	if err := d.Set("instances", instances); err != nil {
		return fmt.Errorf("Error setting instances: %s", err)
	}

	// Set the top-level fields
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}

	// Set the data source ID
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances", project, zone))

	return nil
}
