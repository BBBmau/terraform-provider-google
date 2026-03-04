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
			tpgresource.HandleListError(ctx, req, &diags, push, stream, "Config Error", fmt.Sprintf("Error setting project: %s", err))
			return
		}

		servicesList, err := BatchRequestReadServices(project, tempData, r.Client)
		if err != nil {
			tpgresource.HandleListError(ctx, req, &diags, push, stream, "API Error", err.Error())
			return
		}

		for serviceName := range servicesList.(map[string]struct{}) {
			if err := tempData.Set("project", project); err != nil {
				tpgresource.HandleListError(ctx, req, &diags, push, stream, "Config Error", fmt.Sprintf("Error setting project: %s", err))
				return
			}
			if err := tempData.Set("service", serviceName); err != nil {
				tpgresource.HandleListError(ctx, req, &diags, push, stream, "Config Error", fmt.Sprintf("Error setting service: %s", err))
				return
			}

			tempData.SetId(fmt.Sprintf("%s/%s", project, serviceName))

			identity, err := tempData.Identity()
			if err != nil {
				tpgresource.HandleListError(ctx, req, &diags, push, stream, "Identity Error", fmt.Sprintf("Error getting identity: %s", err))
				return
			}
			if err := identity.Set("project", project); err != nil {
				tpgresource.HandleListError(ctx, req, &diags, push, stream, "Identity Error", fmt.Sprintf("Error setting project on identity: %s", err))
				return
			}
			if err := identity.Set("service", serviceName); err != nil {
				tpgresource.HandleListError(ctx, req, &diags, push, stream, "Identity Error", fmt.Sprintf("Error setting service on identity: %s", err))
				return
			}

			result := req.NewListResult(ctx)
			tfTypeIdentity, err := tempData.TfTypeIdentityState()
			if err != nil {
				tpgresource.HandleListError(ctx, req, &diags, push, stream, "Schema Error", err.Error())
				return
			}
			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				tpgresource.HandleListError(ctx, req, &diags, push, stream, "Schema Error", "error setting identity")
				return
			}
			if req.IncludeResource {
				tfTypeResource, err := tempData.TfTypeResourceState()
				if err != nil {
					tpgresource.HandleListError(ctx, req, &diags, push, stream, "Schema Error", err.Error())
					return
				}
				if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
					tpgresource.HandleListError(ctx, req, &diags, push, stream, "Schema Error", "error setting resource")
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
