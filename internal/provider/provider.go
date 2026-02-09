// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	resp.Schema = schema.Schema{}
}

func (p *playfabProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *playfabProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *playfabProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
