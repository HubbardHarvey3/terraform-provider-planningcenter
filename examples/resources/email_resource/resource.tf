resource "planningcenter_people" "new_test" {
  first_name         = "Tester"
  last_name          = "McTesterson"
  site_administrator = false
  gender             = "Male"
  birthdate          = "1980-05-22"
}


resource "planningcenter_email" "import" {
  address  = "import@importantnotreal.com"
  primary  = true
  location = "Other"
  relationships = {
    id = planningcenter_people.new_test.id
  }
}
