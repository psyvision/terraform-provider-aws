package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsCognitoUserPoolClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCognitoUserPoolClientCreate,
		Read:   resourceAwsCognitoUserPoolClientRead,
		Update: resourceAwsCognitoUserPoolClientUpdate,
		Delete: resourceAwsCognitoUserPoolClientDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"client_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"generate_secret": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"user_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"explicit_auth_flows": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCognitoUserPoolClientAuthFlows,
				},
			},

			"read_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"write_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"refresh_token_validity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntegerInRange(0, 3650),
			},
		},
	}
}

func resourceAwsCognitoUserPoolClientCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	params := &cognitoidentityprovider.CreateUserPoolClientInput{
		ClientName: aws.String(d.Get("name").(string)),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	if v, ok := d.GetOk("generate_secret"); ok {
		params.GenerateSecret = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("explicit_auth_flows"); ok {
		params.ExplicitAuthFlows = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("read_attributes"); ok {
		params.ReadAttributes = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("write_attributes"); ok {
		params.WriteAttributes = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("refresh_token_validity"); ok {
		params.RefreshTokenValidity = aws.Int64(int64(v.(int)))
	}

	log.Printf("[DEBUG] Creating Cognito User Pool Client: %s", params)

	resp, err := conn.CreateUserPoolClient(params)

	if err != nil {
		return errwrap.Wrapf("Error creating Cognito User Pool Client: {{err}}", err)
	}

	d.SetId(*resp.UserPoolClient.ClientId)
	d.Set("user_pool_id", *resp.UserPoolClient.UserPoolId)
	d.Set("name", *resp.UserPoolClient.ClientName)

	return resourceAwsCognitoUserPoolClientRead(d, meta)
}

func resourceAwsCognitoUserPoolClientRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	params := &cognitoidentityprovider.DescribeUserPoolClientInput{
		ClientId:   aws.String(d.Id()),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	log.Printf("[DEBUG] Reading Cognito User Pool Client: %s", params)

	resp, err := conn.DescribeUserPoolClient(params)

	if err != nil {
		if isAWSErr(err, "ResourceNotFoundException", "") {
			log.Printf("[WARN] Cognito User Pool Client %s is already gone", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(*resp.UserPoolClient.ClientId)
	d.Set("user_pool_id", *resp.UserPoolClient.UserPoolId)
	d.Set("name", *resp.UserPoolClient.ClientName)
	d.Set("explicit_auth_flows", flattenStringList(resp.UserPoolClient.ExplicitAuthFlows))
	d.Set("read_attributes", flattenStringList(resp.UserPoolClient.ReadAttributes))
	d.Set("write_attributes", flattenStringList(resp.UserPoolClient.WriteAttributes))
	d.Set("refresh_token_validity", resp.UserPoolClient.RefreshTokenValidity)
	d.Set("client_secret", resp.UserPoolClient.ClientSecret)

	return nil
}

func resourceAwsCognitoUserPoolClientUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	params := &cognitoidentityprovider.UpdateUserPoolClientInput{
		ClientId:   aws.String(d.Id()),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	if d.HasChange("explicit_auth_flows") {
		params.ExplicitAuthFlows = expandStringList(d.Get("explicit_auth_flows").([]interface{}))
	}

	if d.HasChange("read_attributes") {
		params.ReadAttributes = expandStringList(d.Get("read_attributes").([]interface{}))
	}

	if d.HasChange("write_attributes") {
		params.WriteAttributes = expandStringList(d.Get("write_attributes").([]interface{}))
	}

	if d.HasChange("refresh_token_validity") {
		params.RefreshTokenValidity = aws.Int64(d.Get("refresh_token_validity").(int64))
	}

	log.Printf("[DEBUG] Updating Cognito User Pool Client: %s", params)

	_, err := conn.UpdateUserPoolClient(params)
	if err != nil {
		return errwrap.Wrapf("Error updating Cognito User Pool Client: {{err}}", err)
	}

	return resourceAwsCognitoUserPoolClientRead(d, meta)
}

func resourceAwsCognitoUserPoolClientDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	params := &cognitoidentityprovider.DeleteUserPoolClientInput{
		ClientId:   aws.String(d.Id()),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	log.Printf("[DEBUG] Deleting Cognito User Pool Client: %s", params)

	_, err := conn.DeleteUserPoolClient(params)

	if err != nil {
		return errwrap.Wrapf("Error deleting Cognito User Pool Client: {{err}}", err)
	}

	return nil
}
