package objectsv2

func ResourceNutanixObjectsV2() *schema.Resource{
	id_address_schema := &schema.Schema{
		Type: schema.TypeList,
		Optional: true,
		Computed: true,
		maxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {
					Type: schema.TypeList,
					Optional: true,
					Computed: true,
					maxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type: schema.TypeString,
								Required: true,
							},
							"prefix_length": {
								Type: schema.TypeInt,
								Optional: true,
								Computed: true,
								Default: 32,
							},
						},
					},
				},
				"ipv6": {
					Type: schema.TypeList,
					Optional: true,
					Computed: true,
					maxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type: schema.TypeString,
								Required: true,
							},
							"prefix_length": {
								Type: schema.TypeInt,
								Optional: true,
								Computed: true,
								Default: 32,
							},
						},
					}
				}
			}
		}
	}

	return &schema.Resource{
		CreateContext: ResourceNutanixObjectsV2Create,
		ReadContext: ResourceNutanixObjectsV2Read,
		UpdateContext: ResourceNutanixObjectsV2Update,
		DeleteContext: ResourceNutanixObjectsV2Delete,
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				maxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"owner_reference_id": {
							Type: schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"owner_user_name": {
							Type: schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"project_reference_id": {
							Type: schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"project_name": {
							Type: schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"category_ids": {
							Type: schema.TypeList,
							Optional: true,
							Computed: true,
						}
					}
				}
			},
			"name": {
				Type: schema.TypeString,
				Required: true,
			},
			"description": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deployment_version": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"num_worker_nodes": {
				Type: schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cluster_ext_id": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			}
			"storage_network_reference": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storageNetwork_vip": {
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				maxItems: 1,
				Elem: id_address_schema
			},
			"storage_network_dns_ip": {
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				maxItems: 1,
				Elem: id_address_schema
			},
			"public_network_reference": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_network_ips": {
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				maxItems: 1,
				Elem: id_address_schema
			},
			"total_capacity_gib": {
				Type: schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{"DEPLOYING_OBJECT_STORE", "OBJECT_STORE_DEPLOYMENT_FAILED", "DELETING_OBJECT_STORE", "OBJECT_STORE_OPERATION_FAILED", "UNDEPLOYED_OBJECT_STORE", "OBJECT_STORE_OPERATION_PENDING", "OBJECT_STORE_AVAILABLE", "OBJECT_STORE_CERT_CREATION_FAILED", "CREATING_OBJECT_STORE_CERT", "OBJECT_STORE_DELETION_FAILED"}, false),
			},
		}
	}
}