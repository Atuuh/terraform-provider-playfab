// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-playfab/internal/playfab"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &playfabProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &playfabProvider{
			version: version,
		}
	}
}

type playfabProvider struct {
	version string
}

func (p *playfabProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "playfab"
	resp.Version = p.version
}

func (p *playfabProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"title_id": schema.StringAttribute{
				Required: true,
			},
			"secret_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

type playfabProviderModel struct {
	TitleId   types.String `tfsdk:"title_id"`
	SecretKey types.String `tfsdk:"secret_key"`
}

func (p *playfabProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config playfabProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.TitleId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("title_id"),
			"Unknown TitleId",
			"The provider cannot create the PlayFab API client as there is an unknown configuration value for the PlayFab Title Id,",
		)
	}

	if config.SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Unknown Playfab Secret Key",
			"The provider cannot create the PlayFab API client as there is an unknown configuration value for the PlayFab Secret Key. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	titleId := os.Getenv("PLAYFAB_TITLEID")
	secretKey := os.Getenv("PLAYFAB_SECRETKEY")

	if !config.TitleId.IsNull() {
		titleId = config.TitleId.ValueString()
	}

	if !config.SecretKey.IsNull() {
		secretKey = config.TitleId.ValueString()
	}

	if titleId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("title_id"),
			"Missing PlayFab TitleId",
			"The provider cannot create the PlayFab API client as there is a missing or empty value for the PlayFab TitleId. "+
				"Set the title id value in the configuration or use the PLAYFAB_TITLEID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if secretKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Missing PlayFab Secret Key",
			"The provider cannot create the PlayFab API client as there is a missing or empty value for the PlayFab Secret Key. "+
				"Set the secret key value in the configuration or use the PLAYFAB_SECRETKEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := playfab.NewClient(&titleId, &secretKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create PlayFab API Client",
			"An unexpected error occurred when creating the PlayFab API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"PlayFab Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *playfabProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCloudScriptDataSource,
	}
}

func (p *playfabProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCloudScriptResource,
	}
}
