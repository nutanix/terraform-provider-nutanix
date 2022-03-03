package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceFoundationUploadHypervisor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFoundationUploadHypervisorCreate,
		ReadContext:   resourceFoundationUploadHypervisorRead,
		DeleteContext: resourceFoundationUploadHypervisorDelete,
		Schema: map[string]*schema.Schema{
			"installer_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"md5sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"in_whitelist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceFoundationUploadHypervisorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conn := meta.(*Client).FoundationClientAPI
	// Prepare request
	request := &foundation.UploadHypervisorInput{}

	filename, ok := d.GetOk("filename")
	if ok {
		request.Filename = *utils.StringPtr(filename.(string))
	}
	installer, ok := d.GetOk("installer_type")
	if ok {
		request.Installer_type = *utils.StringPtr(installer.(string))
	}

	resp, err := conn.FileManagement.UploadHypervisor(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("md5sum", resp.Md5sum)
	d.Set("name", resp.Name)
	d.Set("in_whitelist", resp.In_Whitelist)
	d.SetId(resource.UniqueId())
	return nil
}

func resourceFoundationUploadHypervisorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceFoundationUploadHypervisorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationClientAPI

	id := d.Id()

	log.Printf("[Debug] Destroying the hypervisor with the ID %s", d.Id())

	if err := conn.FileManagement.DeleteHypervisorAOS(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
