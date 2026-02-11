package provider

import (
	"context"
	"fmt"
	"terraform-provider-playfab/internal/playfab"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &cloudScriptResource{}
	_ resource.ResourceWithConfigure = &cloudScriptResource{}
)

func NewCloudScriptResource() resource.Resource {
	return &cloudScriptResource{}
}

type cloudScriptResource struct {
	client *playfab.Client
}

func (r *cloudScriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
}

func (r *cloudScriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			"trigger_type": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

type cloudScriptResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Url         types.String `tfsdk:"url"`
	TriggerType types.String `tfsdk:"trigger_type"`
}

func (r *cloudScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudScriptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	function := playfab.Function{
		Name:        plan.Name.ValueString(),
		Address:     plan.Url.ValueString(),
		TriggerType: plan.TriggerType.ValueString(),
	}

	err := r.client.CreateCloudScriptFunction(&function)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Cloud Script Function",
			"Could not create Cloud Script Function, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cloudScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudScriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "reading cloudscript resource", map[string]interface{}{"name": state.Name.ValueString()})

	function, err := r.client.GetCloudScriptFunction(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading PlayFab Cloud Script Function",
			"Could not read PlayFab Cloud Script Function "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Url = types.StringValue(function.Address)
	state.TriggerType = types.StringValue(function.TriggerType)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *cloudScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *cloudScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *cloudScriptResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*playfab.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *playfab.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
