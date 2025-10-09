package fc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	fc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixFCOnboardNodes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixFCOnboardNodesCreate,
		ReadContext:   resourceNutanixFCOnboardNodesRead,
		UpdateContext: resourceNutanixFCOnboardNodesUpdate,
		DeleteContext: resourceNutanixFCOnboardNodesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ImageMiniTimeout),
		},
		Schema: map[string]*schema.Schema{
			"node_serial": {
				Type:     schema.TypeString,
				Required: true,
			},
			"block_serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"imaged_node_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixFCOnboardNodesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).FoundationCentral
	resp, err := conn.Service.GetImagedNode(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil || resp.ImagedNodeUUID == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("node_serial", resp.NodeSerial); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("block_serial", resp.BlockSerial); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("imaged_node_uuid", resp.ImagedNodeUUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("model", resp.Model); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_state", resp.NodeState); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_type", resp.NodeType); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixFCOnboardNodesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	serial := d.Get("node_serial")

	// Get client connection
	conn := meta.(*conns.Client).FoundationCentral
	hwManagers, err := conn.Service.ListHardwareManagers(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, hwManager := range hwManagers.HardwareManagers {
		hwManagerNodes, err := conn.Service.ListHardwareManagerNodes(ctx, *hwManager.HardwareManagerUUID)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, node := range hwManagerNodes.Nodes {
			if *node.NodeSerial == serial {
				var req fc.CreateOnboardNodeInput
				err := copier.Copy(&req, node)
				if err != nil {
					return diag.FromErr(err)
				}
				req.EntityID = node.NodeSerial
				req.EntityType = utils.StringPtr("intersight_nodes_to_be_onboarded")

				resp, err := conn.Service.CreateOnboardNode(ctx, &req)
				if err != nil {
					return diag.FromErr(err)
				}
				if resp.ImagedNodeUUID == nil {
					return diag.Errorf("returned node uuid is empty")
				}

				d.SetId(*resp.ImagedNodeUUID)
				return resourceNutanixFCOnboardNodesRead(ctx, d, meta)
			}
		}
	}

	return diag.Errorf("Node not found with serial %s", serial)
}

func resourceNutanixFCOnboardNodesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixFCOnboardNodesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).FoundationCentral
	log.Printf("[DEBUG] Deleting onboarded node: %s, %s", d.Get("node_serial").(string), d.Id())
	err := conn.Service.DeleteOnboardNode(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error while Deleting Onboarded Node: UUID(%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
