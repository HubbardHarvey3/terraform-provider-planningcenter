// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	client "github.com/HubbardHarvey3/terraform-planningcenter-client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PeopleResource{}
var _ resource.ResourceWithImportState = &PeopleResource{}

func NewPeopleResource() resource.Resource {
	return &PeopleResource{}
}

// PeopleResource defines the resource implementation.
type PeopleResource struct {
	client *client.PC_Client
}

// PeopleResourceModel describes the resource data model.
type PeopleResourceModel struct {
	Gender             types.String `tfsdk:"gender"`
	ID                 types.String `tfsdk:"id"`
	Site_Administrator types.Bool   `tfsdk:"site_administrator"`
	First_Name         types.String `tfsdk:"first_name"`
	Last_Name          types.String `tfsdk:"last_name"`
	Birthdate          types.String `tfsdk:"birthdate"`
}

func (r *PeopleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_people"
}

func (r *PeopleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "People resource",

		Attributes: map[string]schema.Attribute{
			"gender": schema.StringAttribute{
				MarkdownDescription: "Gender of the person",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Person's ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_administrator": schema.BoolAttribute{
				MarkdownDescription: "Bool that determines if the person has rights as a site administrator",
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Optional: true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First Name of the person",
				Required:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last Name of the person",
				Required:            true,
			},
			"birthdate": schema.StringAttribute{
				MarkdownDescription: "Birth date of the person.  Formatted as YYYY-MM-DD",
				Sensitive:           true,
				Default:             nil,
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (r *PeopleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.PC_Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PeopleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PeopleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map the Plan/Config to the PeopleRootResource type to send to PC
	var responseData client.PeopleRoot
	responseData.Data.Attributes.LastName = data.Last_Name.ValueString()
	responseData.Data.Attributes.FirstName = data.First_Name.ValueString()
	responseData.Data.Attributes.SiteAdministrator = data.Site_Administrator.ValueBool()
	responseData.Data.Attributes.Gender = data.Gender.ValueString()
	responseData.Data.Attributes.Birthdate = data.Birthdate.ValueString()
	responseData.Data.ID = data.ID.ValueString()

	body, err := client.CreatePeople(r.client, r.client.AppID, r.client.Token, &responseData)
	if err != nil {
		resp.Diagnostics.AddError("Error during CreatePeople", fmt.Sprintf("Error : %v", err))
		return
	}

	var jsonBody client.PeopleRoot
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	data.Gender = types.StringValue(jsonBody.Data.Attributes.Gender)
	data.ID = types.StringValue(jsonBody.Data.ID)
	data.Site_Administrator = types.BoolValue(jsonBody.Data.Attributes.SiteAdministrator)
	data.First_Name = types.StringValue(jsonBody.Data.Attributes.FirstName)
	data.Birthdate = types.StringValue(jsonBody.Data.Attributes.Birthdate)
	data.Last_Name = types.StringValue(jsonBody.Data.Attributes.LastName)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PeopleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PeopleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	//Fetch the data

	jsonBody, err := client.GetPeople(r.client, r.client.AppID, r.client.Token, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error during GetPeople", fmt.Sprintf("Error : %v", err))
		return
	}

	// Overwrite the fetched data to the state
	data.Gender = types.StringValue(jsonBody.Data.Attributes.Gender)
	data.ID = types.StringValue(jsonBody.Data.ID)
	data.Site_Administrator = types.BoolValue(jsonBody.Data.Attributes.SiteAdministrator)
	data.First_Name = types.StringValue(jsonBody.Data.Attributes.FirstName)
	data.Birthdate = types.StringValue(jsonBody.Data.Attributes.Birthdate)
	data.Last_Name = types.StringValue(jsonBody.Data.Attributes.LastName)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *PeopleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PeopleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the Plan/Config to the PeopleRootResource type to send to PC
	var responseData client.PeopleRoot
	responseData.Data.Attributes.LastName = data.Last_Name.ValueString()
	responseData.Data.Attributes.FirstName = data.First_Name.ValueString()
	responseData.Data.Attributes.SiteAdministrator = data.Site_Administrator.ValueBool()
	responseData.Data.Attributes.Gender = data.Gender.ValueString()
	responseData.Data.Attributes.Birthdate = data.Birthdate.ValueString()
	responseData.Data.ID = data.ID.ValueString()

	body, err := client.UpdatePeople(r.client, r.client.AppID, r.client.Token, data.ID.ValueString(), &responseData)

	//convert json back into struct
	var jsonBody client.PeopleRoot
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	data.Gender = types.StringValue(jsonBody.Data.Attributes.Gender)
	data.ID = types.StringValue(jsonBody.Data.ID)
	data.Site_Administrator = types.BoolValue(jsonBody.Data.Attributes.SiteAdministrator)
	data.First_Name = types.StringValue(jsonBody.Data.Attributes.FirstName)
	data.Last_Name = types.StringValue(jsonBody.Data.Attributes.LastName)
	data.Birthdate = types.StringValue(jsonBody.Data.Attributes.Birthdate)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PeopleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PeopleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	client.DeletePeople(r.client, r.client.AppID, r.client.Token, data.ID.ValueString())

}

func (r *PeopleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
