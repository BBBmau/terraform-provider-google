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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/monitoring/MonitoredProject.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package monitoring

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceMonitoringMonitoredProjectNameDiffSuppressFunc(k, old, new string, d tpgresource.TerraformResourceDataChange) bool {
	// Don't suppress if values are empty strings
	if old == "" || new == "" {
		return false
	}

	oldShort := tpgresource.GetResourceNameFromSelfLink(old)
	newShort := tpgresource.GetResourceNameFromSelfLink(new)

	// Suppress if short names are equal
	if oldShort == newShort {
		return true
	}

	_, isOldNumErr := tpgresource.StringToFixed64(oldShort)
	isOldNumber := isOldNumErr == nil
	_, isNewNumErr := tpgresource.StringToFixed64(newShort)
	isNewNumber := isNewNumErr == nil

	// Suppress if comparing a project number to project id
	return isOldNumber != isNewNumber
}

func resourceMonitoringMonitoredProjectNameDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return ResourceMonitoringMonitoredProjectNameDiffSuppressFunc(k, old, new, d)
}

func ResourceMonitoringMonitoredProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringMonitoredProjectCreate,
		Read:   resourceMonitoringMonitoredProjectRead,
		Delete: resourceMonitoringMonitoredProjectDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringMonitoredProjectImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		SchemaVersion: 1,

		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceMonitoringMonitoredProjectResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceMonitoringMonitoredProjectUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"metrics_scope": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `Required. The resource name of the existing Metrics Scope that will monitor this project. Example: locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}`,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: resourceMonitoringMonitoredProjectNameDiffSuppress,
				Description:      `Immutable. The resource name of the 'MonitoredProject'. On input, the resource name includes the scoping project ID and monitored project ID. On output, it contains the equivalent project numbers. Example: 'locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}/projects/{MONITORED_PROJECT_ID_OR_NUMBER}'`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The time when this 'MonitoredProject' was created.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceMonitoringMonitoredProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandMonitoringMonitoredProjectName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	obj, err = resourceMonitoringMonitoredProjectEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}/projects")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new MonitoredProject: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "POST",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutCreate),
		Headers:              headers,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringPermissionError},
	})
	if err != nil {
		return fmt.Errorf("Error creating MonitoredProject: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = MonitoringOperationWaitTime(
		config, res, "Creating MonitoredProject", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create MonitoredProject: %s", err)
	}

	log.Printf("[DEBUG] Finished creating MonitoredProject %q: %#v", d.Id(), res)

	return resourceMonitoringMonitoredProjectRead(d, meta)
}

func resourceMonitoringMonitoredProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	name := d.Get("name").(string)
	name = tpgresource.GetResourceNameFromSelfLink(name)
	d.Set("name", name)
	metricsScope := d.Get("metrics_scope").(string)
	metricsScope = tpgresource.GetResourceNameFromSelfLink(metricsScope)
	d.Set("metrics_scope", metricsScope)
	url, err = tpgresource.ReplaceVars(d, config, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}")
	if err != nil {
		return err
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Headers:              headers,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringPermissionError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("MonitoringMonitoredProject %q", d.Id()))
	}

	res, err = resourceMonitoringMonitoredProjectDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing MonitoringMonitoredProject because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("name", flattenMonitoringMonitoredProjectName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading MonitoredProject: %s", err)
	}
	if err := d.Set("create_time", flattenMonitoringMonitoredProjectCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading MonitoredProject: %s", err)
	}

	return nil
}

func resourceMonitoringMonitoredProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{MonitoringBasePath}}v1/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting MonitoredProject %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "DELETE",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutDelete),
		Headers:              headers,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringPermissionError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "MonitoredProject")
	}

	err = MonitoringOperationWaitTime(
		config, res, "Deleting MonitoredProject", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting MonitoredProject %q: %#v", d.Id(), res)
	return nil
}

func resourceMonitoringMonitoredProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	name := d.Get("name").(string)
	name = tpgresource.GetResourceNameFromSelfLink(name)
	d.Set("name", name)
	metricsScope := d.Get("metrics_scope").(string)
	metricsScope = tpgresource.GetResourceNameFromSelfLink(metricsScope)
	d.Set("metrics_scope", metricsScope)
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"locations/global/metricsScopes/(?P<metrics_scope>[^/]+)/projects/(?P<name>[^/]+)",
		"v1/locations/global/metricsScopes/(?P<metrics_scope>[^/]+)/projects/(?P<name>[^/]+)",
		"(?P<metrics_scope>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenMonitoringMonitoredProjectName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenMonitoringMonitoredProjectCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandMonitoringMonitoredProjectName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func resourceMonitoringMonitoredProjectEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Encoded monitored project name: %s", name)
	name = tpgresource.GetResourceNameFromSelfLink(name)
	log.Printf("[DEBUG] Encoded monitored project resource name: %s", name)
	d.Set("name", name)
	metricsScope := d.Get("metrics_scope").(string)
	log.Printf("[DEBUG] Encoded monitored project metricsScope: %s", metricsScope)
	metricsScope = tpgresource.GetResourceNameFromSelfLink(metricsScope)
	log.Printf("[DEBUG] Encoded monitored project metricsScope resource name: %s", metricsScope)
	d.Set("metrics_scope", metricsScope)
	obj["name"] = fmt.Sprintf("locations/global/metricsScopes/%s/projects/%s", metricsScope, name)
	return obj, nil
}

func resourceMonitoringMonitoredProjectDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	// terraform resource config
	config := meta.(*transport_tpg.Config)

	// The API returns all monitored projects
	monitoredProjects, _ := res["monitoredProjects"].([]interface{})

	// Convert configured terraform monitored_project resource name to a ProjectNumber
	expectedProject, configProjectErr := config.NewResourceManagerClient(config.UserAgent).Projects.Get(d.Get("name").(string)).Do()
	if configProjectErr != nil {
		return nil, configProjectErr
	}
	expectedProjectNumber := strconv.FormatInt(expectedProject.ProjectNumber, 10)

	log.Printf("[DEBUG] Scanning for ProjectNumber: %s.", expectedProjectNumber)

	// Iterate through the list of monitoredProjects to make sure one matches the configured monitored_project
	for _, monitoredProjectRaw := range monitoredProjects {
		if monitoredProjectRaw == nil {
			continue
		}
		monitoredProject := monitoredProjectRaw.(map[string]interface{})

		// MonitoredProject names have the format locations/global/metricsScopes/[metricScopeProjectNumber]/projects/[projectNumber]
		monitoredProjectName := monitoredProject["name"]

		// `res` contains the MonitoredProjects of the relevant metrics scope
		log.Printf("[DEBUG] Matching ProjectNumbers: %s to %s.", expectedProjectNumber, monitoredProjectName)
		if strings.HasSuffix(monitoredProjectName.(string), fmt.Sprintf("/%s", expectedProjectNumber)) {
			// Match found - set response object name to match
			res["name"] = monitoredProjectName
			log.Printf("[DEBUG] Matched ProjectNumbers: %s and %s.", expectedProjectNumber, monitoredProjectName)
			return res, nil
		}
	}
	log.Printf("[DEBUG] MonitoringMonitoredProject couldn't be matched.")
	return nil, nil
}

func resourceMonitoringMonitoredProjectResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metrics_scope": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `Required. The resource name of the existing Metrics Scope that will monitor this project. Example: locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}`,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `Immutable. The resource name of the 'MonitoredProject'. On input, the resource name includes the scoping project ID and monitored project ID. On output, it contains the equivalent project numbers. Example: 'locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}/projects/{MONITORED_PROJECT_ID_OR_NUMBER}'`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The time when this 'MonitoredProject' was created.`,
			},
		},
		UseJSONNumber: true,
	}
}

func ResourceMonitoringMonitoredProjectUpgradeV0(_ context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)

	rawState["id"] = strings.TrimPrefix(rawState["id"].(string), "v1/")

	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}
