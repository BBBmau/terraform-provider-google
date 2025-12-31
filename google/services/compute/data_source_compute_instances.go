package compute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	project := d.Get("project").(string)
	if project == "" {
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("Error getting project: %s", err)
		}
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
	}

	zone := d.Get("zone").(string)
	if zone == "" {
		zone, err := tpgresource.GetZone(d, config)
		if err != nil {
			return fmt.Errorf("Error getting zone: %s", err)
		}
		if err := d.Set("zone", zone); err != nil {
			return fmt.Errorf("Error setting zone: %s", err)
		}
	}

	instances := make([]map[string]interface{}, 0)
	err := ListInstances(context.Background(), d, config, func(item interface{}) error {
		// Convert the map item to compute.Instance struct
		itemMap := item.(map[string]interface{})
		itemJSON, err := json.Marshal(itemMap)
		if err != nil {
			return fmt.Errorf("Error marshaling instance item: %s", err)
		}

		var instance compute.Instance
		if err := json.Unmarshal(itemJSON, &instance); err != nil {
			return fmt.Errorf("Error unmarshaling instance: %s", err)
		}

		name := instance.Name

		tempData := ResourceComputeInstance().Data(nil)
		tempData.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, name))
		tempData.Set("project", project)
		tempData.Set("zone", zone)
		tempData.Set("name", name)

		// Flatten the instance into the temporary ResourceData using the instance from the LIST call
		// compute instance is a very niche case, so normally we would handle it with the following definition:
		// flattenComputeInstance(item.(map[string]interface{}), tempData, config)
		// worth noting that res is already a map[string]interface{} type
		if err := flattenComputeInstanceData(tempData, config, &instance); err != nil {
			return fmt.Errorf("Error flattening instance: %s", err)
		}

		// Extract all values from the temporary ResourceData into a map
		// Use the data source schema to determine which fields to extract
		dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstance().Schema)
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
	// Set the data source ID
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances", project, zone))

	return nil
}

// this would be an example of how it would be implemented for a MMv1 plural datasource
// func dataSourceGoogleComputeInstancesRead(d *schema.ResourceData, meta interface{}) error {
// 	config := meta.(*transport_tpg.Config)

// 	project := d.Get("project").(string)
// 	if project == "" {
// 		project, err := tpgresource.GetProject(d, config)
// 		if err != nil {
// 			return fmt.Errorf("Error getting project: %s", err)
// 		}
// 		if err := d.Set("project", project); err != nil {
// 			return fmt.Errorf("Error setting project: %s", err)
// 		}
// 	}

// 	zone := d.Get("zone").(string)
// 	if zone == "" {
// 		zone, err := tpgresource.GetZone(d, config)
// 		if err != nil {
// 			return fmt.Errorf("Error getting zone: %s", err)
// 		}
// 		if err := d.Set("zone", zone); err != nil {
// 			return fmt.Errorf("Error setting zone: %s", err)
// 		}
// 	}

// 	instances := make([]map[string]interface{}, 0)
// 	err := ListInstances(context.Background(), d, config, func(item map[string]interface{}) error {
// 		// item is already a map[string]interface{} type
//		tempData := ResourceComputeInstance().Data(nil)
//		tempData.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, item["name"].(string)))
//		tempData.Set("project", project)
//		tempData.Set("zone", zone)
//		tempData.Set("name", item["name"].(string))
//		if err := flattenComputeInstance(item,tempData, config); err != nil {
//			return fmt.Errorf("Error flattening instance: %s", err)
//		}
// 		instances = append(instances, tempData.State().Attributes)
// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("instances", instances); err != nil {
// 		return fmt.Errorf("Error setting instances: %s", err)
// 	}
// 	// Set the data source ID
// 	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances", project, zone))

// 	return nil
// }
