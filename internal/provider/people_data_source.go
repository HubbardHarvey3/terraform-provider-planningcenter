// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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

type Root struct {
	Links interface{} `json:"links"`
	Data  Person      `json:"data"`
}
type Person struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	AccountingAdministrator bool        `json:"accounting_administrator"`
	Anniversary             interface{} `json:"anniversary"`
	Avatar                  string      `json:"avatar"`
	Birthdate               string      `json:"birthdate"`
	CanCreateForms          bool        `json:"can_create_forms"`
	CanEmailLists           bool        `json:"can_email_lists"`
	Child                   bool        `json:"child"`
	CreatedAt               time.Time   `json:"created_at"`
	DemographicAvatarURL    string      `json:"demographic_avatar_url"`
	DirectoryStatus         string      `json:"directory_status"`
	FirstName               string      `json:"first_name"`
	Gender                  string      `json:"gender"`
	GivenName               interface{} `json:"given_name"`
	Grade                   interface{} `json:"grade"`
	GraduationYear          interface{} `json:"graduation_year"`
	InactivatedAt           interface{} `json:"inactivated_at"`
	LastName                string      `json:"last_name"`
	MedicalNotes            interface{} `json:"medical_notes"`
	Membership              string      `json:"membership"`
	MiddleName              interface{} `json:"middle_name"`
	Name                    string      `json:"name"`
	Nickname                interface{} `json:"nickname"`
	PassedBackgroundCheck   bool        `json:"passed_background_check"`
	PeoplePermissions       string      `json:"people_permissions"`
	RemoteID                interface{} `json:"remote_id"`
	SchoolType              interface{} `json:"school_type"`
	SiteAdministrator       bool        `json:"site_administrator"`
	Status                  string      `json:"status"`
	UpdatedAt               time.Time   `json:"updated_at"`
}

func NewPeopleDataSource() datasource.DataSource {
	return &PeopleDataSource{}
}

// PeopleDataSource defines the data source implementation.
type PeopleDataSource struct {
	client *PC_Client
}

// PeopleDataSourceModel describes the data source data model.
type PeopleDataSourceModel struct {
	Gender             types.String `tfsdk:"gender"`
	Id                 string       `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Site_Administrator types.Bool   `tfsdk:"site_administrator"`
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
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the person",
				Optional:            true,
			},
			"site_administrator": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (d *PeopleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*PC_Client)

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

	//Fetch the data
	app_id := os.Getenv("PC_APP_ID")
	secret_token := os.Getenv("PC_SECRET_TOKEN")
	endpoint := "https://api.planningcenteronline.com/people/v2/people/" + data.Id
	request, err := http.NewRequest("GET", endpoint, nil)

	request.SetBasicAuth(app_id, secret_token)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	response, err := d.client.Do(request)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	var jsonBody Root
	//var jsonBody map[string]map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Print(err)
	}

	data.Name = types.StringValue(jsonBody.Data.Attributes.Name)
	data.Gender = types.StringValue(jsonBody.Data.Attributes.Gender)
	data.Site_Administrator = types.BoolValue(jsonBody.Data.Attributes.SiteAdministrator)
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
