package licensingv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/agreements"
	import2 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/error"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixLicensesV2 returns the schema.Resource for the licensing v2
func ResourceNutanixLicensesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixLicensesV2Create,
		ReadContext:   ResourceNutanixLicensesV2Read,
		UpdateContext: ResourceNutanixLicensesV2Update,
		DeleteContext: ResourceNutanixLicensesV2Delete,

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// Create license
func ResourceNutanixLicensesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	userName := d.Get("user_name").(string)
	loginID := d.Get("login_id").(string)
	jobTitle := d.Get("job_title").(string)
	companyName := d.Get("company_name").(string)

	license := import1.EndUser{
		UserName:    utils.StringPtr(userName),
		LoginId:     utils.StringPtr(loginID),
		JobTitle:    utils.StringPtr(jobTitle),
		CompanyName: utils.StringPtr(companyName),
	}

	resp, err := conn.EndUserLicenseAgreementAPIInstance.AddUser(&license)
	if err != nil {
		return diag.FromErr(err)
	}

	respData := resp.Data.GetValue().([]import2.AppMessage)

	aJSON, _ := json.MarshalIndent(respData, "", "  ")
	tflog.Debug(ctx, string(aJSON))

	d.SetId(utils.GenUUID())

	return nil
}

// Read license
func ResourceNutanixLicensesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func ResourceNutanixLicensesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixLicensesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
