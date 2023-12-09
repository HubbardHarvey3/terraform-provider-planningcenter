// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-planningcenter/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EmailResource{}
var _ resource.ResourceWithImportState = &EmailResource{}

func NewEmailResource() resource.Resource {
	return &EmailResource{}
}

// EmailResource defines the resource implementation.
type EmailResource struct {
	client *client.PC_Client
}

// EmailResourceModel describes the resource data model.
type EmailResourceModel struct {
	ID            types.String                   `tfsdk:"id"`
	Primary       types.Bool                     `tfsdk:"primary"`
	Location      types.String                   `tfsdk:"location"`
	Address       types.String                   `tfsdk:"address"`
	Relationships EmailResourceRelationshipModel `tfsdk:"relationships"`
}

type EmailResourceRelationshipModel struct {
	PeopleID types.String `tfsdk:"id"`
}

func (r *EmailResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email"
}

func (r *EmailResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "People resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Email's ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"primary": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Computed: true,
				Optional: true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location of the Email, eg. work or personal",
				Required:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The email address",
				Required:            true,
			},
			"relationships": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "People ID the email is associated with",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}
func (r *EmailResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EmailResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map the Plan/Config to the RootResource type to send to PC
	var responseData client.EmailRoot
	responseData.Data.Attributes.Address = data.Address.ValueString()
	responseData.Data.Attributes.Location = data.Location.ValueString()
  peopleID := data.Relationships.PeopleID.ValueString()
	responseData.Data.Attributes.Primary = data.Primary.ValueBool()

	body := client.CreateEmail(r.client, r.client.AppID, r.client.Token, peopleID, &responseData)

	var jsonBody client.EmailRoot
	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	data.Address = types.StringValue(jsonBody.Data.Attributes.Address)
	data.ID = types.StringValue(jsonBody.Data.ID)
	data.Location = types.StringValue(jsonBody.Data.Attributes.Location)
	data.Primary = types.BoolValue(jsonBody.Data.Attributes.Primary)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EmailResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	//Fetch the data

	jsonBody := client.GetEmail(r.client, r.client.AppID, r.client.Token, data.ID.ValueString())

	// Overwrite the fetched data to the state
	data.Address = types.StringValue(jsonBody.Data.Attributes.Address)
	data.ID = types.StringValue(jsonBody.Data.ID)
	data.Location = types.StringValue(jsonBody.Data.Attributes.Location)
	data.Primary = types.BoolValue(jsonBody.Data.Attributes.Primary)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *EmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EmailResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the Plan/Config to the RootResource type to send to PC
	var responseData client.EmailRoot
	responseData.Data.Attributes.Address = data.Address.ValueString()
	responseData.Data.Attributes.Location = data.Location.ValueString()
	responseData.Data.Attributes.Primary = data.Primary.ValueBool()

	body := client.UpdateEmail(r.client, r.client.AppID, r.client.Token, data.ID.ValueString(), &responseData)

	//convert json back into struct
	var jsonBody client.EmailRoot
	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	data.Address = types.StringValue(jsonBody.Data.Attributes.Address)
	data.ID = types.StringValue(jsonBody.Data.ID)
	data.Location = types.StringValue(jsonBody.Data.Attributes.Location)
	data.Primary = types.BoolValue(jsonBody.Data.Attributes.Primary)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EmailResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

  fmt.Println("DELETE")
  fmt.Println("DELETE")
  fmt.Println(data)

	if resp.Diagnostics.HasError() {
		return
	}
	client.DeleteEmail(r.client, r.client.AppID, r.client.Token, data.ID.ValueString())

}

func (r *EmailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
