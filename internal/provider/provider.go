// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"planningcenter/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure PlanningCenterProvider satisfies various provider interfaces.
var _ provider.Provider = &PlanningCenterProvider{}

// PlanningCenterProvider defines the provider implementation.
type PlanningCenterProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PlanningCenterProviderModel describes the provider data model.
type PlanningCenterProviderModel struct {
	Endpoint    types.String `tfsdk:"endpoint"`
	AppId       types.String `tfsdk:"app_id"`
	SecretToken types.String `tfsdk:"secret_token"`
}

func (p *PlanningCenterProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "planningcenter"
	resp.Version = p.version
}

func (p *PlanningCenterProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "App Id provided by Planning Center after requesting Token",
				Optional:            true,
				Sensitive:           true,
			},
			"secret_token": schema.StringAttribute{
				MarkdownDescription: "Secret Token provided by Planning Center after requesting a Token",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *PlanningCenterProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PlanningCenterProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	app_id := os.Getenv("PC_APP_ID")
	secret_token := os.Getenv("PC_SECRET_TOKEN")

	if app_id == "" {
		resp.Diagnostics.AddAttributeError(path.Root("app_id"), "Missing the App ID Environment Variable", "The environment variable 'PC_APP_ID' is missing and the provider can't communicate with Planning Center")
	}
	if secret_token == "" {
		resp.Diagnostics.AddAttributeError(path.Root("secret_token"), "Missing the Secret Token Environment Variable", "The environment variable 'PC_SECRET_TOKEN' is missing and the provider can't communicate with Planning Center")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Example client configuration for data sources and resources
	client := client.NewPCClient(app_id, secret_token, client.HostURL)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PlanningCenterProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPeopleResource,
		NewEmailResource,
	}
}

func (p *PlanningCenterProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPeopleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PlanningCenterProvider{
			version: version,
		}
	}
}
