package multidomainv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	multidomainConfig "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/management"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// schemaLocalDomainConnectionStatus is a fully Computed block shared by the
// nutanix_local_domain_v2 resource and datasource (see HARD RULE 8 — sharing an
// identical, fully Computed nested block is allowed).
func schemaLocalDomainConnectionStatus() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_statuses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// schemaLocalDomainRegistrationStatus is a fully Computed block shared by the
// nutanix_local_domain_v2 resource and datasource.
func schemaLocalDomainRegistrationStatus() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"days_to_expiry": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"error_stage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_provider_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// schemaLocalDomainLocation is a fully Computed block shared by the
// nutanix_local_domain_v2 resource and datasource.
func schemaLocalDomainLocation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"coordinates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latitude": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"longitude": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenLocalDomainConnectionStatus(status *management.LocalDomainConnectionStatus) []map[string]interface{} {
	if status == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})
	if status.Status != nil {
		m["status"] = status.Status.GetName()
	}
	m["service_statuses"] = flattenLocalDomainServiceStatuses(status.ServiceStatuses)
	return []map[string]interface{}{m}
}

func flattenLocalDomainServiceStatuses(statuses []management.ServiceStatus) []map[string]interface{} {
	if len(statuses) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(statuses))
	for _, s := range statuses {
		m := make(map[string]interface{})
		if s.ServiceName != nil {
			m["service_name"] = s.ServiceName.GetName()
		}
		if s.Status != nil {
			m["status"] = s.Status.GetName()
		}
		out = append(out, m)
	}
	return out
}

func flattenLocalDomainRegistrationStatus(status *management.LocalDomainRegistrationStatus) []map[string]interface{} {
	if status == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})
	if status.DaysToExpiry != nil {
		m["days_to_expiry"] = utils.IntValue(status.DaysToExpiry)
	}
	if status.ErrorStage != nil {
		m["error_stage"] = status.ErrorStage.GetName()
	}
	if status.IdentityProviderExtId != nil {
		m["identity_provider_ext_id"] = utils.StringValue(status.IdentityProviderExtId)
	}
	if status.State != nil {
		m["state"] = status.State.GetName()
	}
	return []map[string]interface{}{m}
}

func flattenLocalDomainLocation(location *multidomainConfig.Location) []map[string]interface{} {
	if location == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})
	m["ext_id"] = utils.StringValue(location.ExtId)
	m["name"] = utils.StringValue(location.Name)
	m["country"] = utils.StringValue(location.Country)
	m["state"] = utils.StringValue(location.State)
	m["coordinates"] = flattenLocalDomainGeographicCoordinates(location.Coordinates)
	return []map[string]interface{}{m}
}

func flattenLocalDomainGeographicCoordinates(coords *multidomainConfig.GeographicCoordinates) []map[string]interface{} {
	if coords == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})
	if coords.Latitude != nil {
		m["latitude"] = utils.Float64Value(coords.Latitude)
	}
	if coords.Longitude != nil {
		m["longitude"] = utils.Float64Value(coords.Longitude)
	}
	return []map[string]interface{}{m}
}
