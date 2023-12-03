package client

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
	Name                    string      `json:"name"`
	MedicalNotes            interface{} `json:"medical_notes"`
	Membership              string      `json:"membership"`
	MiddleName              interface{} `json:"middle_name"`
	Nickname                interface{} `json:"nickname"`
	PeoplePermissions       string      `json:"people_permissions"`
	RemoteID                interface{} `json:"remote_id"`
	SiteAdministrator       bool        `json:"site_administrator"`
	Status                  string      `json:"status"`
}
