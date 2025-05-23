// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/netapp/VolumeQuotaRule.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package netapp

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceNetappVolumeQuotaRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetappVolumeQuotaRuleCreate,
		Read:   resourceNetappVolumeQuotaRuleRead,
		Update: resourceNetappVolumeQuotaRuleUpdate,
		Delete: resourceNetappVolumeQuotaRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetappVolumeQuotaRuleImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"disk_limit_mib": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum allowed capacity in MiB.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource name of the quotaRule.`,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidateEnum([]string{"INDIVIDUAL_USER_QUOTA", "INDIVIDUAL_GROUP_QUOTA", "DEFAULT_USER_QUOTA", "DEFAULT_GROUP_QUOTA"}),
				Description:  `Types of Quota Rule. Possible values: ["INDIVIDUAL_USER_QUOTA", "INDIVIDUAL_GROUP_QUOTA", "DEFAULT_USER_QUOTA", "DEFAULT_GROUP_QUOTA"]`,
			},
			"volume_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Name of the volume to create the quotaRule in.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Description for the quota rule.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Labels as key value pairs of the quota rule. Example: '{ "owner": "Bob", "department": "finance", "purpose": "testing" }'.


**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Loction of the quotaRule. QuotaRules are child resources of volumes and live in the same location.`,
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The quota rule applies to the specified user or group.
Valid targets for volumes with NFS protocol enabled:
  - UNIX UID for individual user quota
  - UNIX GID for individual group quota
Valid targets for volumes with SMB protocol enabled:
  - Windows SID for individual user quota
Leave empty for default quotas`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Create time of the quota rule. A timestamp in RFC3339 UTC "Zulu" format. Examples: "2023-06-22T09:13:01.617Z".`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The state of the quota rule. Possible Values : [STATE_UNSPECIFIED, CREATING, UPDATING, READY, DELETING, ERROR]`,
			},
			"state_details": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `State details of the quota rule`,
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceNetappVolumeQuotaRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	targetProp, err := expandNetappVolumeQuotaRuleTarget(d.Get("target"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target"); !tpgresource.IsEmptyValue(reflect.ValueOf(targetProp)) && (ok || !reflect.DeepEqual(v, targetProp)) {
		obj["target"] = targetProp
	}
	typeProp, err := expandNetappVolumeQuotaRuleType(d.Get("type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("type"); !tpgresource.IsEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	diskLimitMibProp, err := expandNetappVolumeQuotaRuleDiskLimitMib(d.Get("disk_limit_mib"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disk_limit_mib"); !tpgresource.IsEmptyValue(reflect.ValueOf(diskLimitMibProp)) && (ok || !reflect.DeepEqual(v, diskLimitMibProp)) {
		obj["diskLimitMib"] = diskLimitMibProp
	}
	descriptionProp, err := expandNetappVolumeQuotaRuleDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	labelsProp, err := expandNetappVolumeQuotaRuleEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetappBasePath}}projects/{{project}}/locations/{{location}}/volumes/{{volume_name}}/quotaRules?quotaRuleId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new VolumeQuotaRule: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for VolumeQuotaRule: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating VolumeQuotaRule: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/volumes/{{volume_name}}/quotaRules/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = NetappOperationWaitTime(
		config, res, project, "Creating VolumeQuotaRule", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create VolumeQuotaRule: %s", err)
	}

	log.Printf("[DEBUG] Finished creating VolumeQuotaRule %q: %#v", d.Id(), res)

	return resourceNetappVolumeQuotaRuleRead(d, meta)
}

func resourceNetappVolumeQuotaRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetappBasePath}}projects/{{project}}/locations/{{location}}/volumes/{{volume_name}}/quotaRules/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for VolumeQuotaRule: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("NetappVolumeQuotaRule %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}

	if err := d.Set("target", flattenNetappVolumeQuotaRuleTarget(res["target"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("type", flattenNetappVolumeQuotaRuleType(res["type"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("disk_limit_mib", flattenNetappVolumeQuotaRuleDiskLimitMib(res["diskLimitMib"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("state", flattenNetappVolumeQuotaRuleState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("state_details", flattenNetappVolumeQuotaRuleStateDetails(res["stateDetails"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("create_time", flattenNetappVolumeQuotaRuleCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("description", flattenNetappVolumeQuotaRuleDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("labels", flattenNetappVolumeQuotaRuleLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("terraform_labels", flattenNetappVolumeQuotaRuleTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}
	if err := d.Set("effective_labels", flattenNetappVolumeQuotaRuleEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading VolumeQuotaRule: %s", err)
	}

	return nil
}

func resourceNetappVolumeQuotaRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for VolumeQuotaRule: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	targetProp, err := expandNetappVolumeQuotaRuleTarget(d.Get("target"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, targetProp)) {
		obj["target"] = targetProp
	}
	typeProp, err := expandNetappVolumeQuotaRuleType(d.Get("type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("type"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	diskLimitMibProp, err := expandNetappVolumeQuotaRuleDiskLimitMib(d.Get("disk_limit_mib"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disk_limit_mib"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, diskLimitMibProp)) {
		obj["diskLimitMib"] = diskLimitMibProp
	}
	descriptionProp, err := expandNetappVolumeQuotaRuleDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	labelsProp, err := expandNetappVolumeQuotaRuleEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetappBasePath}}projects/{{project}}/locations/{{location}}/volumes/{{volume_name}}/quotaRules/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating VolumeQuotaRule %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("target") {
		updateMask = append(updateMask, "target")
	}

	if d.HasChange("type") {
		updateMask = append(updateMask, "type")
	}

	if d.HasChange("disk_limit_mib") {
		updateMask = append(updateMask, "diskLimitMib")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("effective_labels") {
		updateMask = append(updateMask, "labels")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// if updateMask is empty we are not updating anything so skip the post
	if len(updateMask) > 0 {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
			Headers:   headers,
		})

		if err != nil {
			return fmt.Errorf("Error updating VolumeQuotaRule %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating VolumeQuotaRule %q: %#v", d.Id(), res)
		}

		err = NetappOperationWaitTime(
			config, res, project, "Updating VolumeQuotaRule", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceNetappVolumeQuotaRuleRead(d, meta)
}

func resourceNetappVolumeQuotaRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for VolumeQuotaRule: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{NetappBasePath}}projects/{{project}}/locations/{{location}}/volumes/{{volume_name}}/quotaRules/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting VolumeQuotaRule %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "VolumeQuotaRule")
	}

	err = NetappOperationWaitTime(
		config, res, project, "Deleting VolumeQuotaRule", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting VolumeQuotaRule %q: %#v", d.Id(), res)
	return nil
}

func resourceNetappVolumeQuotaRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/volumes/(?P<volume_name>[^/]+)/quotaRules/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<volume_name>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<volume_name>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/volumes/{{volume_name}}/quotaRules/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNetappVolumeQuotaRuleTarget(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetappVolumeQuotaRuleType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetappVolumeQuotaRuleDiskLimitMib(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenNetappVolumeQuotaRuleState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetappVolumeQuotaRuleStateDetails(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetappVolumeQuotaRuleCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetappVolumeQuotaRuleDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetappVolumeQuotaRuleLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenNetappVolumeQuotaRuleTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("terraform_labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenNetappVolumeQuotaRuleEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNetappVolumeQuotaRuleTarget(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetappVolumeQuotaRuleType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetappVolumeQuotaRuleDiskLimitMib(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetappVolumeQuotaRuleDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetappVolumeQuotaRuleEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
