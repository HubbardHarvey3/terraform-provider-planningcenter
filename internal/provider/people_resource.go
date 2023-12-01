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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
)

//Setup Structs to marshall json into
type RootResource struct {
	Links interface{} `json:"links"`
	Data  PersonResource      `json:"data"`
}
type PersonResource struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes AttributesResource `json:"attributes"`
}

type AttributesResource struct {
	AccountingAdministrator   bool      `json:"accounting_administrator"`
	Anniversary               interface{} `json:"anniversary"`
	Avatar                    string    `json:"avatar"`
	Birthdate                 string    `json:"birthdate"`
	Child                     bool      `json:"child"`
	FirstName                 types.String    `tfsdk:"first_name"`
	Gender                    types.String    `tfsdk:"gender"`
	GivenName                 interface{} `json:"given_name"`
	Grade                     interface{} `json:"grade"`
	GraduationYear            interface{} `json:"graduation_year"`
	InactivatedAt             interface{} `json:"inactivated_at"`
	LastName                  types.String    `tfsdk:"last_name"`
	MedicalNotes              interface{} `json:"medical_notes"`
	Membership                string    `json:"membership"`
	MiddleName                interface{} `json:"middle_name"`
	Nickname                  interface{} `json:"nickname"`
	PeoplePermissions         string    `json:"people_permissions"`
	RemoteID                  interface{} `json:"remote_id"`
	SiteAdministrator         bool      `json:"site_administrator"`
	Status                    string    `json:"status"`
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PeopleResource{}
var _ resource.ResourceWithImportState = &PeopleResource{}

func NewPeopleResource() resource.Resource {
	return &PeopleResource{}
}

// PeopleResource defines the resource implementation.
type PeopleResource struct {
	client *PC_Client
}

// PeopleResourceModel describes the resource data model.
type PeopleResourceModel struct {
	Gender             types.String `tfsdk:"gender"`
	Id                 basetypes.StringValue `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Site_Administrator bool   `tfsdk:"site_administrator"`
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
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the person",
				Optional:            true,
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

	client, ok := req.ProviderData.(*PC_Client)

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

	goBody := RootResource{
		Data: PersonResource{
			Type: "Person",
			Attributes: AttributesResource{
				FirstName:         data.First_Name,
				LastName:          data.Last_Name,
				SiteAdministrator: data.Site_Administrator,
				Gender:            data.Gender,
			},
		},
	}

	// Convert struct to JSON
	jsonData, err := json.MarshalIndent(goBody, "", " ")
  fmt.Println(jsonData)
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
	response, err := r.client.Do(request)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	fmt.Println(response)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	//	tflog.Trace(ctx, "created a resource")



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
	app_id := os.Getenv("PC_APP_ID")
	secret_token := os.Getenv("PC_SECRET_TOKEN")
	endpoint := "https://api.planningcenteronline.com/people/v2/people/" + data.Id.ValueString()
	request, err := http.NewRequest("GET", endpoint, nil)
	request.SetBasicAuth(app_id, secret_token)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	response, err := r.client.Do(request)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	var jsonBody RootResource
	//var jsonBody map[string]map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	data.Gender = jsonBody.Data.Attributes.Gender
	data.Site_Administrator = jsonBody.Data.Attributes.SiteAdministrator
  data.First_Name = jsonBody.Data.Attributes.FirstName
  data.Last_Name = jsonBody.Data.Attributes.LastName

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *PeopleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
