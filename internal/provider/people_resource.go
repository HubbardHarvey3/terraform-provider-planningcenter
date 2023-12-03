// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"terraform-provider-planningcenter/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
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
	Id                 types.String `tfsdk:"id"`
	Site_Administrator types.Bool   `tfsdk:"site_administrator"`
	First_Name         types.String `tfsdk:"first_name"`
	Last_Name          types.String `tfsdk:"last_name"`
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
			},
			"site_administrator": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Computed: true,
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

	//Fetch the data
	app_id := os.Getenv("PC_APP_ID")
	secret_token := os.Getenv("PC_SECRET_TOKEN")
	endpoint := "https://api.planningcenteronline.com/people/v2/people/"

	// Map the Plan/Config to the RootResource type to send to PC
	var responseData client.Root
	responseData.Data.Attributes.LastName = data.Last_Name.ValueString()
	responseData.Data.Attributes.FirstName = data.First_Name.ValueString()
	responseData.Data.Attributes.SiteAdministrator = data.Site_Administrator.ValueBool()
	responseData.Data.Attributes.Gender = data.Gender.ValueString()

	// Convert struct to JSON
	jsonData, err := json.Marshal(&responseData)


	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Create a request with the JSON data
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type to application/json
	request.Header.Set("Content-Type", "application/json")

	// Make the request
	request.SetBasicAuth(app_id, secret_token)
	response, err := r.client.Client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	var jsonBody client.Root
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	data.Gender = types.StringValue(jsonBody.Data.Attributes.Gender)
	data.Id = types.StringValue(jsonBody.Data.ID)
	data.Site_Administrator = types.BoolValue(jsonBody.Data.Attributes.SiteAdministrator)
	data.First_Name = types.StringValue(jsonBody.Data.Attributes.FirstName)
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
	if data.Id.ValueString() == "" {
		fmt.Println("_______________Nothing to read for Nil ID__________________")
		return
	} else {
		//Fetch the data
		app_id := os.Getenv("PC_APP_ID")
		secret_token := os.Getenv("PC_SECRET_TOKEN")

    jsonBody := client.GetPeople(r.client, app_id, secret_token, data.Id.ValueString())

		// Overwrite the fetched data to the state
		data.Gender = types.StringValue(jsonBody.Data.Attributes.Gender)
		data.Id = types.StringValue(jsonBody.Data.ID)
		data.Site_Administrator = types.BoolValue(jsonBody.Data.Attributes.SiteAdministrator)
		data.First_Name = types.StringValue(jsonBody.Data.Attributes.FirstName)
		data.Last_Name = types.StringValue(jsonBody.Data.Attributes.LastName)
		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *PeopleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PeopleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

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
  client.DeletePeople(r.client, r.client.AppID,r.client.Token, data.Id.ValueString())

}

func (r *PeopleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
