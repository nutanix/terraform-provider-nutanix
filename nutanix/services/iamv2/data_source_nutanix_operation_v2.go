package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixOperationV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixOperationV4Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_operation_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"associated_endpoint_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixOperationV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id")

	resp, err := conn.OperationsAPIInstance.GetOperationById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching image placement : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Operation)

	if err := d.Set("display_name", getResp.DisplayName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_type", getResp.EntityType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("operation_type", flattenOperationType(getResp.OperationType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_name", getResp.ClientName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("related_operation_list", utils.StringSlice(getResp.RelatedOperationList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("associated_endpoint_list", flattenAssociatedEndpointList(getResp.AssociatedEndpointList)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		t := getResp.CreatedTime
		if err := d.Set("created_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdatedTime != nil {
		t := getResp.LastUpdatedTime
		if err := d.Set("last_updated_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenOperationType(pr *import1.OperationType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == import1.OperationType(two) {
			return "INTERNAL"
		}
		if *pr == import1.OperationType(three) {
			return "SYSTEM_DEFINED_ONLY"
		}
		if *pr == import1.OperationType(four) {
			return "EXTERNAL"
		}
	}
	return "UNKNOWN"
}

func flattenAssociatedEndpointList(pr []import1.AssociatedEndpoint) []map[string]interface{} {
	if len(pr) > 0 {
		endpoints := make([]map[string]interface{}, len(pr))
		for _, v := range pr {
			endpoint := make(map[string]interface{})

			endpoint["api_version"] = flattenAPIVersion(v.ApiVersion)
			endpoint["endpoint_url"] = v.EndpointUrl
			endpoint["http_method"] = flattenHTTPMethod(v.HttpMethod)

			endpoints = append(endpoints, endpoint)
		}
		return endpoints
	}
	return nil
}

func flattenAPIVersion(pr *import1.ApiVersion) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import1.ApiVersion(two) {
			return "V3"
		}
		if *pr == import1.ApiVersion(three) {
			return "V4"
		}
	}
	return "UNKNOWN"
}

func flattenHTTPMethod(pr *import1.HttpMethod) string {
	if pr != nil {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		if *pr == import1.HttpMethod(two) {
			return "HTTPMETHOD_POST"
		}
		if *pr == import1.HttpMethod(three) {
			return "HTTPMETHOD_GET"
		}
		if *pr == import1.HttpMethod(four) {
			return "HTTPMETHOD_PUT"
		}
		if *pr == import1.HttpMethod(five) {
			return "HTTPMETHOD_PATCH"
		}
		if *pr == import1.HttpMethod(six) {
			return "HTTPMETHOD_DELETE"
		}
	}
	return "UNKNOWN"
}
