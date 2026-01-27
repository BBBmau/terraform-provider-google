package compute

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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
		Blocks: map[string]listschema.Block{
			"filter": customListFiltersBlock(),
		},
	}
}

type ComputeInstanceListModel struct {
	Project types.String            `tfsdk:"project"`
	Zone    types.String            `tfsdk:"zone"`
	Filter  []customListFilterModel `tfsdk:"filter"`
}

func customListFiltersBlock() listschema.ListNestedBlock {
	return listschema.ListNestedBlock{
		NestedObject: listschema.NestedBlockObject{
			Attributes: map[string]listschema.Attribute{
				"name": listschema.StringAttribute{
					Required: true,
				},
				"values": listschema.ListAttribute{
					Required:    true,
					ElementType: types.StringType,
				},
			},
		},
	}
}

type customListFilterModel struct {
	Name   types.String        `tfsdk:"name"`
	Values basetypes.ListValue `tfsdk:"values"`
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

	var zone string
	if !data.Zone.IsNull() && !data.Zone.IsUnknown() {
		zone = data.Zone.ValueString()
	}
	if zone == "" {
		zone = r.Client.Zone
	}

	filterString := ""
	if len(data.Filter) > 0 {
		for i, filter := range data.Filter {
			values := make([]string, 0)
			diags := filter.Values.ElementsAs(ctx, &values, false)
			if diags.HasError() {
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			// values represents the operator and value used for filtering
			filterString += fmt.Sprintf("(%s %s \"%s\")", filter.Name.ValueString(), values[0], values[1])
			if i < len(data.Filter)-1 {
				filterString += " AND "
			}
		}
	}

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListInstances(r.Client, filterString, func(rd *schema.ResourceData) error {
			result := req.NewListResult(ctx)

			// flatten using the instance from the LIST call
			identity, err := rd.Identity()
			if err != nil {
				return fmt.Errorf("Error getting identity: %s", err)
			}
			err = identity.Set("name", rd.Get("name").(string))
			if err != nil {
				return fmt.Errorf("Error setting name: %s", err)
			}
			err = identity.Set("zone", zone)
			if err != nil {
				return fmt.Errorf("Error setting zone: %s", err)
			}
			err = identity.Set("project", project)
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
			if req.IncludeResource {
				tfTypeResource, err := rd.TfTypeResourceState()
				if err != nil {
					return err
				}
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

func ListInstances(config *transport_tpg.Config, filter string, callback func(rd *schema.ResourceData) error) error {
	url, err := tpgresource.ReplaceVars(ResourceComputeInstance().Data(&terraform.InstanceState{}), config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances")
	if err != nil {
		return err
	}

	opts := ListCallOptions{
		Config:    config,
		TempData:  ResourceComputeInstance().Data(&terraform.InstanceState{}),
		Url:       url,
		Filter:    filter,
		Flattener: flattenComputeInstance,
		Callback:  callback,
	}

	return ListCall(opts)
}

type ListCallOptions struct {
	Config    *transport_tpg.Config
	TempData  *schema.ResourceData
	Url       string
	ItemName  string
	Filter    string
	Flattener func(item interface{}, d *schema.ResourceData, config *transport_tpg.Config) error
	Callback  func(rd *schema.ResourceData) error
}

func ListCall(opts ListCallOptions) error {
	// Set default ItemName if not provided
	if opts.ItemName == "" {
		opts.ItemName = "items"
	}

	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(opts.Url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(opts.TempData, opts.Config); err == nil {
		billingProject = bp
	}

	userAgent, err := tpgresource.GenerateUserAgentString(opts.TempData, opts.Config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	if opts.Filter != "" {
		params["filter"] = opts.Filter
	}

	for {
		// Depending on previous iterations, params might contain a pageToken param
		url, err := transport_tpg.AddQueryParams(opts.Url, params)
		if err != nil {
			return err
		}

		headers := make(http.Header)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    opts.Config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Headers:   headers,
			// ErrorRetryPredicates used to allow retrying if rate limits are hit when requesting multiple pages in a row
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, opts.TempData, fmt.Sprintf("%s %q", opts.ItemName, opts.TempData.Id()))
		}

		if v, ok := res[opts.ItemName].([]interface{}); ok {
			for _, item := range v {

				err = opts.Flattener(item, opts.TempData, opts.Config)
				if err != nil {
					return fmt.Errorf("Error flattening instance: %s", err)
				}
				err = opts.Callback(opts.TempData)
				if err != nil {
					return err
				}
			}
		}
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
