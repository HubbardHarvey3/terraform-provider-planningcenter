---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "planningcenter_email Resource - terraform-provider-planningcenter"
subcategory: ""
description: |-
  Email resource
---

# planningcenter_email (Resource)

Email resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) The email address
- `location` (String) Location of the Email, eg. work or personal
- `relationships` (Attributes) (see [below for nested schema](#nestedatt--relationships))

### Optional

- `primary` (Boolean) Whether or not the email is the person's primary or not

### Read-Only

- `id` (String) Email's ID

<a id="nestedatt--relationships"></a>
### Nested Schema for `relationships`

Required:

- `id` (String) People ID the email is associated with
