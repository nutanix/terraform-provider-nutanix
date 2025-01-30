package selfservice

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func ResourceNutanixCalmEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmEndpointCreate,
		ReadContext:   resourceNutanixCalmEndpointRead,
		UpdateContext: resourceNutanixCalmEndpointUpdate,
		DeleteContext: resourceNutanixCalmEndpointDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cred_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cred_password": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNutanixCalmEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conn := meta.(*conns.Client).Calm
	name := d.Get("name").(string)
	desc := d.Get("description").(string)
	ipAddress := d.Get("ip_address").(string)
	credUsername := d.Get("cred_username").(string)
	credPassword := d.Get("cred_password").(string)

	metadata := createMetadata(meta)
	spec := createSpec(name, desc, ipAddress, credUsername, credPassword)

	endpointInput := &calm.EndpointCreateInput{}

	endpointInput.Spec = spec
	endpointInput.Metadata = metadata
	endpointInput.APIVersion = "3.0"

	createResp, err := conn.Service.CreateEndpoint(ctx, endpointInput)
	if err != nil {
		return diag.FromErr(err)
	}

	epUUID := createResp.Metadata["uuid"].(string)

	d.SetId(epUUID)

	return nil
}

func resourceNutanixCalmEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixCalmEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixCalmEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func createMetadata(meta interface{}) map[string]interface{} {
	metadata := map[string]interface{}{}
	metadata["kind"] = "endpoint"
	newUuid, err := uuid.GenerateUUID()

	if err != nil {
		fmt.Println("Error while creating uuid.", err)
	}

	metadata["uuid"] = newUuid

	apiConn := meta.(*conns.Client).API
	resp, err := apiConn.V3.ListAllProject("")

	if err != nil {
		fmt.Println("Error while fetching projects.", err)
	}

	projRef := map[string]interface{}{}
	projRef["name"] = resp.Entities[0].Status.Name
	projRef["kind"] = "project"
	projRef["uuid"] = resp.Entities[0].Status.UUID

	metadata["project_reference"] = projRef

	return metadata
}

func createSpec(name string, desc string, ipAddress string, credUsername string, credPassword string) map[string]interface{} {
	spec := map[string]interface{}{}
	resources := map[string]interface{}{}

	attrs := map[string]interface{}{}

	credDef := map[string]interface{}{}
	credDef["description"] = ""
	credDef["username"] = credUsername
	credDef["type"] = "PASSWORD"
	credDef["cred_class"] = "static"
	nameUUID, _ := uuid.GenerateUUID()
	credUUID, _ := uuid.GenerateUUID()
	credDef["name"] = "endpoint_cred_" + nameUUID[:8]
	credDef["uuid"] = credUUID

	secret := map[string]interface{}{}
	secAttrs := map[string]interface{}{}
	secAttrs["is_secret_modified"] = true
	secret["attrs"] = secAttrs
	secret["value"] = credPassword

	credDef["secret"] = secret

	credDefList := []map[string]interface{}{}
	credDefList = append(credDefList, credDef)

	loginCredRef := map[string]interface{}{}
	loginCredRef["name"] = credDef["name"]
	loginCredRef["uuid"] = credDef["uuid"]
	loginCredRef["kind"] = "app_credential"

	attrs["port"] = "22"
	attrs["values"] = [1]string{ipAddress}
	attrs["credential_definition_list"] = credDefList
	attrs["login_credential_reference"] = loginCredRef

	resources["attrs"] = attrs
	resources["type"] = "Linux"
	resources["value_type"] = "IP"

	spec["name"] = name
	spec["description"] = desc
	spec["resources"] = resources
	return spec
}
