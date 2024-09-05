package ephemeral_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewPasswordEphemeral() ephemeral.EphemeralResource {
	return &passwordEphemeral{}
}

type passwordEphemeral struct{}

func (p *passwordEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

type ephemeralPasswordModel struct {
	Length types.Int64  `tfsdk:"length"`
	Result types.String `tfsdk:"result"`
}

func (p *passwordEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a random password",
		Attributes: map[string]schema.Attribute{
			"length": schema.Int64Attribute{
				Description: "The length of the string desired. The minimum value for length is 1.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"result": schema.StringAttribute{
				Description: "The generated random string.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *passwordEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralPasswordModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result := "test"

	data.Result = types.StringValue(string(result))
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
