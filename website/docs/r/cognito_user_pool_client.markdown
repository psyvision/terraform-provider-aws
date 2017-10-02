layout: "aws"
page_title: "AWS: aws_cognito_user_pool_client"
side_bar_current: "docs-aws-resource-cognito-user-pool-client"
description: |-
  Provides a Cognito User Pool Client resource.

# aws_cognito_user_pool_client

Provides a Cognito User Pool Client resource.

## Example Usage

### Create a basic user pool client client

```hcl
resource "aws_cognito_user_pool_client" "client" {
  name = "client"
}
```

### Create a user pool client with no SRP authentication
```hcl
resource "aws_cognito_user_pool_client" "client" {
  name = "client"

  generate_secret = true
  explicit_auth_flows = ["ADMIN_NO_SRP_AUTH"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user pool.
* `generate_secret` - (Optional) The subject line for verification emails.
* `user_pool_id` - (Required) The user pool the client belongs to.
* `explicit_auth_flows` - (Optional) List of authentication flows (ADMIN_NO_SRP_AUTH, CUSTOM_AUTH_FLOW_ONLY)

## Attribute Reference

The following attributes are exported:

* `id` - The id of the user pool client.
* `client_secret` - The client secret of the user pool client.