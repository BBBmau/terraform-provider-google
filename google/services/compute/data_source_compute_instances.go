package compute

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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
	err := ListInstances(context.Background(), config, func(rd *schema.ResourceData) error {
		// Extract all values from the temporary ResourceData into a map
		// Use the data source schema to determine which fields to extract
		dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstance().Schema)
		instanceMap := make(map[string]interface{})
		for key := range dsSchema {
			if val, ok := rd.GetOk(key); ok {
				// Only include non-nil values
				if val != nil {
					instanceMap[key] = val
				}
			}
		}
		log.Printf("rd.State().Attributes: %+v", rd.State().Attributes)
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

//	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstance().Schema)
//	instanceMap := make(map[string]interface{})
//	for key := range dsSchema {
//		if val, ok := tempData.GetOk(key); ok {
//			// Only include non-nil values
//			if val != nil {
//				instanceMap[key] = val
//			}
//		}
//	}
// 		instances = append(instances, instanceMap)
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
