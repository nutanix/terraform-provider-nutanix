package selfservice

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func ResourceNutanixCalmBlueprintCreate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmBlueprintCreate,
		ReadContext:   resourceNutanixCalmBlueprintRead,
		UpdateContext: resourceNutanixCalmBlueprintUpdate,
		DeleteContext: resourceNutanixCalmBlueprintDelete,
		Schema: map[string]*schema.Schema{
			"bp_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nic_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"credentials": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixCalmBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm
	bp_name := d.Get("bp_name").(string)
	// project_uuid := d.Get("project_uuid").(string)
	filePath := "/Users/abhinav.bansal1/Downloads/create_bp.json"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("content::: %+v", string(content))

	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	jsonData["spec"].(map[string]interface{})["name"] = bp_name

	var blueprint calm.CreateBlueprintResponse

	// Unmarshal JSON into the Go struct
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonBytes, &blueprint)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("blueprint::: %+v", blueprint)

	bp, err := conn.Service.CreateBlueprint(ctx, blueprint)
	if err != nil {
		return diag.Errorf("Error creating blueprint: %s", err)
	}
	log.Printf("bp::: %+v", bp)

	// Print spec of bp
	log.Printf("bp.Spec::: %+v", bp.Spec)

	jsonData1, err := json.Marshal(bp.Spec)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("jsonData1::: %+v", string(jsonData1))

	credentialDefinitionList := bp.Spec["resources"].(map[string]interface{})["credential_definition_list"].([]interface{})
	credentialDefinitionList[0].(map[string]interface{})["secret"] = map[string]interface{}{
		"attrs": map[string]interface{}{
			"is_secret_modified": true,
		},
		"value": "nutanix/4u",
	}

	// log credentialDefinitionList and bp.Spec
	log.Printf("credentialDefinitionList::: %+v", credentialDefinitionList)
	log.Printf("bp.Spec_abhi::: %+v", bp.Spec)

	bp_uuid := bp.Metadata["uuid"].(string)
	bp_update, err := conn.Service.UpdateBlueprint(ctx, bp_uuid, *bp)
	if err != nil {
		return diag.Errorf("Error updating blueprint: %s", err)
	}
	log.Printf("bp_update::: %+v", bp_update)
	// Set the ID of the resource in Terraform
	d.SetId(bp_uuid)
	return nil
}

func resourceNutanixCalmBlueprintRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
