package compute

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var _ tpgresource.ListResourceWithRawV5Schemas = &ComputeInstanceListResource{}

type ComputeInstanceListResource struct {
	tpgresource.ListResourceMetadata
}

func NewComputeInstanceListResource() list.ListResource {
	return &ComputeInstanceListResource{}
}

func (r *ComputeInstanceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "google_compute_instance"
}

func (r *ComputeInstanceListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	computeInstance := ResourceComputeInstance()
	resp.ProtoV5Schema = computeInstance.ProtoSchema(ctx)()
	resp.ProtoV5IdentitySchema = computeInstance.ProtoIdentitySchema(ctx)()
}

func (r *ComputeInstanceListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *ComputeInstanceListResource) ListResourceConfigSchema(ctx context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"project": listschema.StringAttribute{
				Optional: true,
			},
			"zone": listschema.StringAttribute{
				Optional: true,
			},
		},
	}
}

type ComputeInstanceListModel struct {
	Project types.String `tfsdk:"project"`
	Zone    types.String `tfsdk:"zone"`
}

func (r *ComputeInstanceListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data ComputeInstanceListModel
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	var project string
	if !data.Project.IsNull() && !data.Project.IsUnknown() {
		project = data.Project.ValueString()
	}
	if project == "" {
		project = r.Client.Project
	}
	r.Client.Project = project

	var zone string
	if !data.Zone.IsNull() && !data.Zone.IsUnknown() {
		zone = data.Zone.ValueString()
	}
	if zone == "" {
		zone = r.Client.Zone
	}

	stream.Results = func(push func(list.ListResult) bool) {
		computeInstanceResource := ResourceComputeInstance()
		rd := computeInstanceResource.Data(&terraform.InstanceState{})
		rd.Set("project", project)
		rd.Set("zone", zone)
		// This is how it would be called in a plural datasource
		// err := ListInstances(ctx, rd, r.Client, func(item interface{}) error {
		// 	instances = append(instances, flattenComputeInstance(item))
		// return nil
		// })
		err := ListInstances(ctx, rd, r.Client, func(item interface{}) error {
			result := req.NewListResult(ctx)
			result.DisplayName = item.(map[string]interface{})["name"].(string)
			rd.Set("name", result.DisplayName)
			// flatten
			if err := flattenComputeInstance(rd, r.Client); err != nil {
				result.Diagnostics.AddError("Error flattening instance: %s", err.Error())
				return err
			}
			identity, err := rd.Identity()
			if err != nil {
				return fmt.Errorf("Error getting identity: %s", err)
			}
			err = identity.Set("name", result.DisplayName)
			if err != nil {
				return fmt.Errorf("Error setting name: %s", err)
			}
			err = identity.Set("zone", zone)
			if err != nil {
				return fmt.Errorf("Error setting zone: %s", err)
			}
			err = identity.Set("project", r.Client.Project)
			if err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				return err
			}
			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				return errors.New("error setting identity")
			}
			tfTypeResource, err := rd.TfTypeResourceState()
			if err != nil {
				return err
			}
			if req.IncludeResource {
				if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
					return errors.New("error setting resource")
				}
			}
			if !push(result) {
				return errors.New("stream closed")
			}
			return nil
		})
		if err != nil {
			diags.AddError("API Error", err.Error())
			result := req.NewListResult(ctx)
			result.Diagnostics = diags
			push(result)
		}
		stream.Results = list.ListResultsStreamDiagnostics(diags)
	}
}

func ListInstances(ctx context.Context, rd *schema.ResourceData, config *transport_tpg.Config, callback func(interface{}) error) error {
	url, err := tpgresource.ReplaceVars(rd, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances")
	if err != nil {
		return err
	}

	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(rd, config); err == nil {
		billingProject = bp
	}

	userAgent, err := tpgresource.GenerateUserAgentString(rd, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	if v, ok := rd.GetOk("filter"); ok {
		params["filter"] = v.(string)
	}

	for {
		// Depending on previous iterations, params might contain a pageToken param
		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}
		log.Printf("url: %s\n", url)

		headers := make(http.Header)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Headers:   headers,
			// ErrorRetryPredicates used to allow retrying if rate limits are hit when requesting multiple pages in a row
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, rd, fmt.Sprintf("Instances %q", rd.Id()))
		}

		// we need to figure out what to do here since this is where we get the response from the LIST API
		// Store info from this page

		// this currently works because we use the callback for every
		if v, ok := res["items"].([]interface{}); ok {
			for _, item := range v {
				// TODO: add a flattener function here
				// no point in adding the flatten as part of the callback function
				// flattener that sets the values by rd.Set

				// NOTE: for our case we need a flattener that flattens the ENTIRE resource, there currently is no existing
				// flattener that puts all the existing flatteners together.

				// to prevent overlap we can work on refactoring it later so we can just focus on the
				// root flattener to being used only by list and plural datasources.
				// rd.Set("items", result.DisplayName)

				// flatten
				err = callback(item)
				if err != nil {
					return err
				}
			}
		}
		// ------------------------------------------------------------------------------------------------

		// Handle pagination for next loop, or break loop
		v, ok := res["nextPageToken"]
		if ok {
			params["pageToken"] = v.(string)
		}
		if !ok {
			break
		}
	}
	return nil
}
