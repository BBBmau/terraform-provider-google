package resourcemanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

var _ tpgresource.ListResourceWithRawV5Schemas = &ProjectServiceListResource{}

type ProjectServiceListResource struct {
	tpgresource.ListResourceMetadata
}

func NewProjectServiceListResource() list.ListResource {
	return &ProjectServiceListResource{}
}

func (r *ProjectServiceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "google_project_service"
}

func (r *ProjectServiceListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	projectService := ResourceGoogleProjectService()
	resp.ProtoV5Schema = projectService.ProtoSchema(ctx)()
	resp.ProtoV5IdentitySchema = projectService.ProtoIdentitySchema(ctx)()
}

func (r *ProjectServiceListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *ProjectServiceListResource) ListResourceConfigSchema(ctx context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"project": listschema.StringAttribute{
				Optional: true,
			},
		},
	}
}

type ProjectServiceListModel struct {
	Project types.String `tfsdk:"project"`
}

func (r *ProjectServiceListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data ProjectServiceListModel
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
	project = tpgresource.GetResourceNameFromSelfLink(project)

	stream.Results = func(push func(list.ListResult) bool) {
		// Use a temporary ResourceData for BatchRequestReadServices (needs it for user agent, billing project, timeout)
		tempData := ResourceGoogleProjectService().Data(&terraform.InstanceState{})
		if err := tempData.Set("project", project); err != nil {
			diags.AddError("Config Error", fmt.Sprintf("Error setting project: %s", err))
			result := req.NewListResult(ctx)
			result.Diagnostics = diags
			push(result)
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}

		servicesRaw, err := BatchRequestReadServices(project, tempData, r.Client)
		if err != nil {
			diags.AddError("API Error", err.Error())
			result := req.NewListResult(ctx)
			result.Diagnostics = diags
			push(result)
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}

		servicesList, ok := servicesRaw.(map[string]struct{})
		if !ok {
			diags.AddError("API Error", "unexpected type from ListCurrentlyEnabledServices")
			result := req.NewListResult(ctx)
			result.Diagnostics = diags
			push(result)
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}

		resourceData := ResourceGoogleProjectService().Data(&terraform.InstanceState{})

		for serviceName := range servicesList {
			if err := resourceData.Set("project", project); err != nil {
				diags.AddError("Config Error", fmt.Sprintf("Error setting project: %s", err))
				result := req.NewListResult(ctx)
				result.Diagnostics = diags
				push(result)
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			if err := resourceData.Set("service", serviceName); err != nil {
				diags.AddError("Config Error", fmt.Sprintf("Error setting service: %s", err))
				result := req.NewListResult(ctx)
				result.Diagnostics = diags
				push(result)
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}

			resourceData.SetId(fmt.Sprintf("%s/%s", project, serviceName))

			identity, err := resourceData.Identity()
			if err != nil {
				diags.AddError("Identity Error", fmt.Sprintf("Error getting identity: %s", err))
				result := req.NewListResult(ctx)
				result.Diagnostics = diags
				push(result)
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			if err := identity.Set("project", project); err != nil {
				diags.AddError("Identity Error", fmt.Sprintf("Error setting project on identity: %s", err))
				result := req.NewListResult(ctx)
				result.Diagnostics = diags
				push(result)
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			if err := identity.Set("service", serviceName); err != nil {
				diags.AddError("Identity Error", fmt.Sprintf("Error setting service on identity: %s", err))
				result := req.NewListResult(ctx)
				result.Diagnostics = diags
				push(result)
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}

			result := req.NewListResult(ctx)
			tfTypeIdentity, err := resourceData.TfTypeIdentityState()
			if err != nil {
				diags.AddError("Schema Error", err.Error())
				result.Diagnostics = diags
				push(result)
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				diags.AddError("Schema Error", "error setting identity")
				result.Diagnostics = diags
				push(result)
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			if req.IncludeResource {
				tfTypeResource, err := resourceData.TfTypeResourceState()
				if err != nil {
					diags.AddError("Schema Error", err.Error())
					result.Diagnostics = diags
					push(result)
					stream.Results = list.ListResultsStreamDiagnostics(diags)
					return
				}
				if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
					diags.AddError("Schema Error", "error setting resource")
					result.Diagnostics = diags
					push(result)
					stream.Results = list.ListResultsStreamDiagnostics(diags)
					return
				}
			}
			if !push(result) {
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
		}

		stream.Results = list.ListResultsStreamDiagnostics(diags)
	}
}
