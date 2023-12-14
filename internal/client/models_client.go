package client

import ()

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
	Child                   bool        `json:"child"`
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
	Nickname                interface{} `json:"nickname"`
	PeoplePermissions       string      `json:"people_permissions"`
	RemoteID                interface{} `json:"remote_id"`
	SiteAdministrator       bool        `json:"site_administrator"`
	Status                  string      `json:"status"`
}

type EmailRootNoRelationship struct {
	Data EmailNoRelationship `json:"data,omitempty"`
}

type EmailNoRelationship struct {
	Type       string          `json:"type"`
	ID         string          `json:"id"`
	Attributes EmailAttributes `json:"attributes"`
}

type EmailRoot struct {
	Data Email `json:"data,omitempty"`
}

type Email struct {
	Type          string             `json:"type"`
	ID            string             `json:"id"`
	Attributes    EmailAttributes    `json:"attributes"`
	Relationships EmailRelationships `json:"relationships,omitempty"`
}

type EmailAttributes struct {
	Address  string `json:"address"`
	Location string `json:"location"`
	Primary  bool   `json:"primary"`
}

type EmailRelationships struct {
	Person EmailPerson `json:"person,omitempty"`
}

type EmailPerson struct {
	Data EmailPersonData `json:"data,omitempty"`
}

type EmailPersonData struct {
	Type string `json:"person,omitempty"`
	ID   string `json:"id,omitempty"`
}
