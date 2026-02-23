// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-playfab/internal/playfab"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &cloudScriptDataSource{}
	_ datasource.DataSourceWithConfigure = &cloudScriptDataSource{}
)

func NewCloudScriptDataSource() datasource.DataSource {
	return &cloudScriptDataSource{}
}

type cloudScriptDataSource struct {
	client *playfab.Client
}

func (d *cloudScriptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_script"
}

func (d *cloudScriptDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"functions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"trigger_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type cloudScriptDataSourceModel struct {
	Functions []functionModel `tfsdk:"functions"`
}

type functionModel struct {
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	TriggerType types.String `tfsdk:"trigger_type"`
}

func (d *cloudScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudScriptDataSourceModel

	functions, err := d.client.GetCloudScriptFunctions()
	ctx = tflog.SetField(ctx, "functions", functions)
	ctx = tflog.SetField(ctx, "f", ctx.Value("title_entity_token"))
	tflog.Info(ctx, "functions datasource read")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PlayFab Cloud Script Functions",
			err.Error(),
		)
		return
	}

	for _, function := range functions {
		functionState := functionModel{
			Name:        types.StringValue(function.Name),
			Address:     types.StringValue(function.Address),
			TriggerType: types.StringValue(function.TriggerType),
		}

		state.Functions = append(state.Functions, functionState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning(fmt.Sprintf("we have %d functions.\n%+v", len(state.Functions), state), "stuff")
}

func (d *cloudScriptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*playfab.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *playfab.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}
