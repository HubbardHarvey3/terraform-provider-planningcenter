// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"planningcenter/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &PeopleDataSource{}
	_ datasource.DataSourceWithConfigure = &PeopleDataSource{}
)

func NewPeopleDataSource() datasource.DataSource {
	return &PeopleDataSource{}
}

// PeopleDataSource defines the data source implementation.
type PeopleDataSource struct {
	client *client.PC_Client
}

// PeopleDataSourceModel describes the data source data model.
type PeopleDataSourceModel struct {
	Gender             types.String `tfsdk:"gender"`
	Id                 types.String `tfsdk:"id"`
	Site_Administrator types.Bool   `tfsdk:"site_administrator"`
	FirstName          types.String `tfsdk:"first_name"`
	LastName           types.String `tfsdk:"last_name"`
}

func (d *PeopleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_people"
}

func (d *PeopleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "People Data Source",

		Attributes: map[string]schema.Attribute{
			"gender": schema.StringAttribute{
				MarkdownDescription: "Gender of the person",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Person's ID",
				Required:            true,
			},
			"site_administrator": schema.BoolAttribute{
				MarkdownDescription: "Bool that determines if the person has rights as a site administrator",
				Computed:            true,
				Optional:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First Name of the person",
				Optional:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last Name of the person",
				Optional:            true,
			},
		},
	}
}

func (d *PeopleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.PC_Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected PC_Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PeopleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PeopleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	app_id := os.Getenv("PC_APP_ID")
	secret_token := os.Getenv("PC_SECRET_TOKEN")

	//Fetch Data from the PC API
	jsonBody := client.GetPeople(d.client, app_id, secret_token, data.Id.ValueString())

	data.Gender = types.StringValue(jsonBody.Data.Attributes.Gender)
	data.Id = types.StringValue(jsonBody.Data.ID)
	data.Site_Administrator = types.BoolValue(jsonBody.Data.Attributes.SiteAdministrator)
	data.FirstName = types.StringValue(jsonBody.Data.Attributes.FirstName)
	data.LastName = types.StringValue(jsonBody.Data.Attributes.LastName)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
