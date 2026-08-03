package networkingv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	networkingConfig "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	networkingBridgesReq "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/bridges"
	networkingVsReq "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/virtualswitches"
	import4 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import5 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVirtualSwitchV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixVirtualSwitchV2Create,
		ReadContext:   resourceNutanixVirtualSwitchV2Read,
		UpdateContext: resourceNutanixVirtualSwitchV2Update,
		DeleteContext: resourceNutanixVirtualSwitchV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV2(),
				},
			},
			"project_ext_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "UUID of the project that owns this entity",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User-visible Virtual Switch name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Input body to configure a Virtual Switch",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is a default Virtual Switch which cannot be deleted",
			},
			"is_quick_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When true, the node is not put in maintenance mode during the Virtual Switch update operation.",
			},
			"mtu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MTU",
			},
			"bond_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE_BACKUP", "BALANCE_SLB", "BALANCE_TCP", "NONE"}, false),
			},
			"clusters": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Cluster configuration list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Reference ExtId for the cluster.",
						},
						"gateway_ip_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The IPv4 address of the host.",
									},
									"prefix_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The prefix length of the network to which this host IPv4 address belongs.",
									},
								},
							},
						},
						"existing_bridge_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Existing bridge name",
						},
						"hosts": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Host configuration array",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Reference to the host",
									},
									"host_nics": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Host NIC array",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"internal_bridge_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Internal bridge name as br0",
									},

									"ip_address": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"prefix_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Prefix length of the IPv4 subnet.",
												},
											},
										},
									},
									"active_uplink": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Host active uplink interface",
									},
									"route_table": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Internal route table number for the routing rules associated with the IP address on this host",
									},
								},
							},
						},
						"vlan_identifier": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "VLAN Identifier for this virtual switch cluster; set to 0 to remove VLAN tagging.",
						},
					},
				},
			},
			"igmp_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_snooping_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable IGMP snooping on this Virtual Switch",
						},
						"snooping_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "IGMP Snooping timeout value in seconds",
						},
						"querier_spec": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_querier_enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Enable IGMP querier on this Virtual Switch",
									},
									"vlan_id_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "VLAN Id list on which IGMP queries must be sent",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
					},
				},
			},
			"owner_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PC", "PE"}, false),
				Description:  "Owner type.",
			},
			"shared_with_projects": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of project UUIDs this virtual switch is shared with.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"has_deployment_error": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When true, virtual switch configuration is not deployed on every node.",
			},
			"has_update_in_progress": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the virtual switch's update is being processed",
			},
			"has_delete_in_progress": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the virtual switch's delete is being processed",
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// detectMigrationFromExistingBridge inspects the `clusters` block in the
// user's config and returns the (bridgeName, clusterRef, isMigration)
// tuple needed by the migrate code path.
//
// Detection rule: ANY entry in `clusters` whose `existing_bridge_name` is a
// non-empty string flips the resource into migrate mode. The migrate API
// only accepts a single cluster reference, so when more than one cluster
// entry sets `existing_bridge_name` we surface a clear error.
//
// Returns isMigration=false (zero values for the rest) when no entry sets
// `existing_bridge_name` -- in that case the standard create path runs.
func detectMigrationFromExistingBridge(d *schema.ResourceData) (bridgeName, clusterRef string, isMigration bool, err error) {
	clustersRaw, ok := d.GetOk("clusters")
	if !ok {
		return "", "", false, nil
	}
	clusters, ok := clustersRaw.([]interface{})
	if !ok || len(clusters) == 0 {
		return "", "", false, nil
	}

	var matched int
	for i, c := range clusters {
		m, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		ebn, _ := m["existing_bridge_name"].(string)
		if ebn == "" {
			continue
		}
		matched++
		if matched > 1 {
			return "", "", true, fmt.Errorf(
				"existing_bridge_name set on more than one clusters[] entry; the migrate API only accepts a single cluster reference -- declare it on clusters[0] only")
		}
		bridgeName = ebn
		clusterRef, _ = m["ext_id"].(string)
		_ = i
	}
	return bridgeName, clusterRef, matched > 0, nil
}

// waitVirtualSwitchTaskAndSetID polls the prism task associated with a
// virtual-switch create/migrate call, extracts the resulting VS UUID from
// the task's EntitiesAffected list, and sets it as the resource ID.
// Shared by the standard-create and migrate code paths.
func waitVirtualSwitchTaskAndSetID(ctx context.Context, d *schema.ResourceData, meta interface{}, taskUUID *string) diag.Diagnostics {
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Virtual Switch (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIDRequest := import5.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIDRequest)
	if err != nil {
		return diag.Errorf("error while fetching Virtual Switch task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVirtualSwitch, "Virtual switch")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return nil
}

func resourceNutanixVirtualSwitchV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bridgeName, clusterRef, isMigration, err := detectMigrationFromExistingBridge(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if isMigration {
		return resourceNutanixVirtualSwitchV2CreateFromBridge(ctx, d, meta, bridgeName, clusterRef)
	}
	return resourceNutanixVirtualSwitchV2StandardCreate(ctx, d, meta)
}

// resourceNutanixVirtualSwitchV2StandardCreate posts to
// POST /api/networking/v4.3/config/virtual-switches with the full
// VirtualSwitch body. Used when no clusters[].existing_bridge_name is set.
func resourceNutanixVirtualSwitchV2StandardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := networkingConfig.VirtualSwitch{}

	if metadata, ok := d.GetOk("metadata"); ok {
		inputSpec.Metadata = expandMetadata(metadata.([]interface{}))
	}
	if ownerType, ok := d.GetOk("owner_type"); ok {
		inputSpec.OwnerType = common.ExpandEnum[networkingConfig.OwnerType](ownerType.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		inputSpec.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		inputSpec.Description = utils.StringPtr(description.(string))
	}
	if bondMode, ok := d.GetOk("bond_mode"); ok {
		inputSpec.BondMode = common.ExpandEnum[networkingConfig.BondModeType](bondMode.(string))
	}
	if clusters, ok := d.GetOk("clusters"); ok {
		inputSpec.Clusters = expandVirtualSwitchClusters(clusters.([]interface{}))
	}
	if mtu, ok := d.GetOk("mtu"); ok {
		inputSpec.Mtu = utils.Int64Ptr(int64(mtu.(int)))
	}
	if igmpSpec, ok := d.GetOk("igmp_spec"); ok {
		inputSpec.IgmpSpec = expandIgmpSpec(igmpSpec.([]interface{}))
	}
	if isQuickMode, ok := d.GetOk("is_quick_mode"); ok {
		inputSpec.IsQuickMode = utils.BoolPtr(isQuickMode.(bool))
	}
	if projectExtID, ok := d.GetOk("project_ext_id"); ok {
		inputSpec.ProjectExtId = utils.StringPtr(projectExtID.(string))
	}

	createReq := networkingVsReq.CreateVirtualSwitchRequest{
		Body: &inputSpec,
	}

	aJSON, _ := json.MarshalIndent(createReq, "", " ")
	log.Printf("[DEBUG] VirtualSwitch create payload: %s", string(aJSON))

	resp, err := conn.VirtualSwitchAPI.CreateVirtualSwitch(ctx, &createReq)
	if err != nil {
		return diag.Errorf("error while creating Virtual Switch: %v", err)
	}

	taskRef := resp.Data.GetValue().(import4.TaskReference)
	if diags := waitVirtualSwitchTaskAndSetID(ctx, d, meta, taskRef.ExtId); diags != nil {
		return diags
	}

	// Share the freshly-created Virtual Switch with each configured project
	// through the dedicated share endpoint. The standard create endpoint
	// does not honor shared_with_projects, so this must be done explicitly.
	if sharedWithProjects, ok := d.GetOk("shared_with_projects"); ok {
		projectIDs := common.ExpandListOfString(sharedWithProjects.([]interface{}))
		if diags := shareVirtualSwitchWithProjects(ctx, meta, d.Id(), projectIDs, d.Timeout(schema.TimeoutCreate)); diags != nil {
			return diags
		}
	}

	return resourceNutanixVirtualSwitchV2Read(ctx, d, meta)
}

// resourceNutanixVirtualSwitchV2CreateFromBridge posts to
// POST /api/networking/v4.3/config/virtual-switches/$actions/migrate
// using the BridgesAPI. The migrate body (config.Bridge) is a STRICT
// subset of the standard create body: only name, description,
// existingBridgeName, clusterReference, projectExtId, and metadata are
// accepted. Everything else (bond_mode, mtu, igmp_spec, host_nics,
// gateway_ip_address, ...) is silently dropped by the API even though the
// schema allows it -- we log a warning in that case so the user has a
// chance to notice. We don't error: it's common for HCL to carry
// standard-create defaults that aren't relevant to migrate.
func resourceNutanixVirtualSwitchV2CreateFromBridge(ctx context.Context, d *schema.ResourceData, meta interface{}, bridgeName, clusterRef string) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	warnIgnoredMigrateFields(d)

	bridge := networkingConfig.Bridge{
		ExistingBridgeName: utils.StringPtr(bridgeName),
	}
	if clusterRef != "" {
		bridge.ClusterReference = utils.StringPtr(clusterRef)
	}
	if name, ok := d.GetOk("name"); ok {
		bridge.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		bridge.Description = utils.StringPtr(description.(string))
	}
	if projectExtID, ok := d.GetOk("project_ext_id"); ok {
		bridge.ProjectExtId = utils.StringPtr(projectExtID.(string))
	}
	if metadata, ok := d.GetOk("metadata"); ok {
		bridge.Metadata = expandMetadata(metadata.([]interface{}))
	}

	migrateReq := networkingBridgesReq.MigrateBridgeRequest{
		Body: &bridge,
	}

	aJSON, _ := json.MarshalIndent(migrateReq, "", " ")
	log.Printf("[DEBUG] VirtualSwitch migrate payload: %s", string(aJSON))

	resp, err := conn.BridgesAPI.MigrateBridge(ctx, &migrateReq)
	if err != nil {
		return diag.Errorf("error while migrating bridge %q to Virtual Switch: %v", bridgeName, err)
	}

	taskRef := resp.Data.GetValue().(import4.TaskReference)
	if diags := waitVirtualSwitchTaskAndSetID(ctx, d, meta, taskRef.ExtId); diags != nil {
		return diags
	}
	return resourceNutanixVirtualSwitchV2Read(ctx, d, meta)
}

// warnIgnoredMigrateFields logs a warning for every attribute the migrate API
// does not accept but the user has set in HCL. Non-fatal -- the migrate API
// will succeed regardless, the warning just makes it obvious why bond_mode etc.
// don't take effect.
//
// The migrate body (config.Bridge) only carries name, description,
// existingBridgeName, clusterReference (from clusters[].ext_id), projectExtId,
// and metadata. Two groups of fields are therefore dropped:
//  1. top-level fields (bond_mode, mtu, igmp_spec, ...), and
//  2. cluster/host-level fields nested under clusters[] (gateway_ip_address,
//     vlan_identifier, and the entire hosts[] block).
func warnIgnoredMigrateFields(d *schema.ResourceData) {
	type ignored struct{ field, hint string }
	candidates := []ignored{
		{"bond_mode", "migrate does not set bondMode -- the resulting VS inherits the source bridge's mode"},
		{"mtu", "migrate does not set mtu -- update the VS after migration if needed"},
		{"igmp_spec", "migrate does not configure igmp_spec -- update the VS after migration if needed"},
		{"is_quick_mode", "migrate does not honor is_quick_mode"},
		{"shared_with_projects", "migrate does not honor shared_with_projects -- set it via update after migration"},
		{"owner_type", "migrate does not honor owner_type"},
	}
	for _, c := range candidates {
		if _, set := d.GetOk(c.field); set {
			log.Printf("[WARN] virtual_switch_v2 migrate: field %q is set but the migrate API ignores it (%s)", c.field, c.hint)
		}
	}

	// Only clusters[].ext_id (-> clusterReference) and
	// clusters[].existing_bridge_name survive into the migrate body. Warn for
	// any other cluster/host-level attribute the user configured.
	clustersRaw, ok := d.GetOk("clusters")
	if !ok {
		return
	}
	clusters, ok := clustersRaw.([]interface{})
	if !ok {
		return
	}
	for i, cRaw := range clusters {
		c, ok := cRaw.(map[string]interface{})
		if !ok {
			continue
		}
		if gw, ok := c["gateway_ip_address"].([]interface{}); ok && len(gw) > 0 {
			log.Printf("[WARN] virtual_switch_v2 migrate: clusters[%d].gateway_ip_address is set but the migrate API ignores it -- update the VS after migration if needed", i)
		}
		if vlan, ok := c["vlan_identifier"].(int); ok && vlan != 0 {
			log.Printf("[WARN] virtual_switch_v2 migrate: clusters[%d].vlan_identifier is set but the migrate API ignores it -- update the VS after migration if needed", i)
		}
		if hosts, ok := c["hosts"].([]interface{}); ok && len(hosts) > 0 {
			log.Printf("[WARN] virtual_switch_v2 migrate: clusters[%d].hosts is set but the migrate API ignores host-level configuration (host_nics, ip_address, active_uplink, route_table); the VS inherits the source bridge's host bindings", i)
		}
	}
}

func resourceNutanixVirtualSwitchV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	getReq := networkingVsReq.GetVirtualSwitchByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.VirtualSwitchAPI.GetVirtualSwitchById(ctx, &getReq)
	if err != nil {
		d.SetId("")
		return nil
	}

	vs := resp.Data.GetValue().(networkingConfig.VirtualSwitch)

	if err := d.Set("ext_id", vs.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", vs.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", vs.Description); err != nil {
		return diag.FromErr(err)
	}
	if vs.BondMode != nil {
		if err := d.Set("bond_mode", vs.BondMode.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	// `existing_bridge_name` is a create-time-only input to the migrate API
	// (see resourceNutanixVirtualSwitchV2CreateFromBridge). The
	// GetVirtualSwitchById response does not echo it back, so the flatten
	// helper writes "" for it. To keep terraform plan stable across
	// refreshes we overlay the prior state value at the matching cluster
	// index. On import (no prior state) the field stays empty -- that's
	// fine because the bridge has already been migrated and the value is
	// only useful at create time.
	flatClusters := flattenVirtualSwitchClusters(vs.Clusters)
	priorClustersRaw, _ := d.GetOk("clusters")
	if priorClusters, ok := priorClustersRaw.([]interface{}); ok {
		for i := range flatClusters {
			if i >= len(priorClusters) {
				break
			}
			priorMap, ok := priorClusters[i].(map[string]interface{})
			if !ok {
				continue
			}
			if prior, ok := priorMap["existing_bridge_name"].(string); ok && prior != "" {
				flatClusters[i]["existing_bridge_name"] = prior
			}
		}
	}
	if err := d.Set("clusters", flatClusters); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mtu", vs.Mtu); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("igmp_spec", flattenIgmpSpec(vs.IgmpSpec)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_quick_mode", vs.IsQuickMode); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_default", vs.IsDefault); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_deployment_error", vs.HasDeploymentError); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_update_in_progress", vs.HasUpdateInProgress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_delete_in_progress", vs.HasDeleteInProgress); err != nil {
		return diag.FromErr(err)
	}
	if vs.OwnerType != nil {
		if err := d.Set("owner_type", vs.OwnerType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("project_ext_id", vs.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("shared_with_projects", vs.SharedWithProjects); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(vs.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(vs.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", vs.TenantId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// virtualSwitchUpdatableFields enumerates every attribute that is part of the
// VirtualSwitch body sent to UpdateVirtualSwitchById. shared_with_projects is
// deliberately excluded: it is reconciled through the dedicated
// share/unshare endpoints, not the main update body. project_ext_id is also
// excluded: although the API accepts it on update, changing the owning project
// of an existing virtual switch is intentionally not supported (the UI client
// behaves the same), so a change to it is rejected with an error instead.
var virtualSwitchUpdatableFields = []string{
	"name",
	"description",
	"bond_mode",
	"clusters",
	"mtu",
	"igmp_spec",
	"is_quick_mode",
}

func resourceNutanixVirtualSwitchV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// The owning project of an existing virtual switch cannot be changed. The
	// API technically supports it, but updating it is intentionally disallowed
	// to match the UI client behaviour.
	if d.HasChange("project_ext_id") {
		return diag.Errorf("project_ext_id cannot be updated: Change of project_ext_id is not supported.")
	}

	// Reconcile any change to the body fields through the standard update
	// endpoint. shared_with_projects is NOT part of this body -- handling it
	// here would conflict with the explicit share/unshare calls below, so it
	// is reconciled separately.
	if d.HasChanges(virtualSwitchUpdatableFields...) {
		getReq := networkingVsReq.GetVirtualSwitchByIdRequest{
			ExtId: utils.StringPtr(d.Id()),
		}
		readResp, err := conn.VirtualSwitchAPI.GetVirtualSwitchById(ctx, &getReq)
		if err != nil {
			return diag.Errorf("error while fetching Virtual Switch for update: %v", err)
		}

		updateSpec := readResp.Data.GetValue().(networkingConfig.VirtualSwitch)

		if d.HasChange("name") {
			updateSpec.Name = utils.StringPtr(d.Get("name").(string))
		}
		if d.HasChange("description") {
			updateSpec.Description = utils.StringPtr(d.Get("description").(string))
		}
		if d.HasChange("bond_mode") {
			updateSpec.BondMode = common.ExpandEnum[networkingConfig.BondModeType](d.Get("bond_mode").(string))
		}
		if d.HasChange("clusters") {
			updateSpec.Clusters = expandVirtualSwitchClusters(common.InterfaceToSlice(d.Get("clusters")))
		}
		if d.HasChange("mtu") {
			updateSpec.Mtu = utils.Int64Ptr(int64(d.Get("mtu").(int)))
		}
		if d.HasChange("igmp_spec") {
			updateSpec.IgmpSpec = expandIgmpSpec(common.InterfaceToSlice(d.Get("igmp_spec")))
		}
		if d.HasChange("is_quick_mode") {
			updateSpec.IsQuickMode = utils.BoolPtr(d.Get("is_quick_mode").(bool))
		}

		etagValue := conn.VirtualSwitchAPI.ApiClient.GetEtag(readResp)
		args := make(map[string]interface{})
		args["If-Match"] = etagValue

		updateReq := networkingVsReq.UpdateVirtualSwitchByIdRequest{
			ExtId: utils.StringPtr(d.Id()),
			Body:  &updateSpec,
		}

		aJSON, _ := json.MarshalIndent(updateReq, "", " ")
		log.Printf("[DEBUG] VirtualSwitch update payload: %s", string(aJSON))

		resp, err := conn.VirtualSwitchAPI.UpdateVirtualSwitchById(ctx, &updateReq, args)
		if err != nil {
			return diag.Errorf("error while updating Virtual Switch: %v", err)
		}

		taskRef := resp.Data.GetValue().(import4.TaskReference)
		if err := waitForVirtualSwitchTask(ctx, meta, taskRef.ExtId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for Virtual Switch (%s) to update: %s", utils.StringValue(taskRef.ExtId), err)
		}
	}

	// Reconcile sharing: compute the projects that were added and removed,
	// then share/unshare each one through the dedicated endpoints.
	if d.HasChange("shared_with_projects") {
		oldRaw, newRaw := d.GetChange("shared_with_projects")
		oldProjects := common.ExpandListOfString(oldRaw.([]interface{}))
		newProjects := common.ExpandListOfString(newRaw.([]interface{}))
		added, removed := diffProjectLists(oldProjects, newProjects)

		if diags := shareVirtualSwitchWithProjects(ctx, meta, d.Id(), added, d.Timeout(schema.TimeoutUpdate)); diags != nil {
			return diags
		}
		if diags := unshareVirtualSwitchFromProjects(ctx, meta, d.Id(), removed, d.Timeout(schema.TimeoutUpdate)); diags != nil {
			return diags
		}
	}

	return resourceNutanixVirtualSwitchV2Read(ctx, d, meta)
}

func resourceNutanixVirtualSwitchV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	deleteReq := networkingVsReq.DeleteVirtualSwitchByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	resp, err := conn.VirtualSwitchAPI.DeleteVirtualSwitchById(ctx, &deleteReq)
	if err != nil {
		return diag.Errorf("error while deleting Virtual Switch: %v", err)
	}

	taskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Virtual Switch (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

// diffProjectLists compares the previous and desired project-UUID lists and
// returns the UUIDs that need to be shared (present in new, absent from old)
// and the UUIDs that need to be unshared (present in old, absent from new).
func diffProjectLists(oldProjects, newProjects []string) (added, removed []string) {
	oldSet := make(map[string]struct{}, len(oldProjects))
	for _, p := range oldProjects {
		oldSet[p] = struct{}{}
	}
	newSet := make(map[string]struct{}, len(newProjects))
	for _, p := range newProjects {
		newSet[p] = struct{}{}
	}

	for _, p := range newProjects {
		if _, ok := oldSet[p]; !ok {
			added = append(added, p)
		}
	}
	for _, p := range oldProjects {
		if _, ok := newSet[p]; !ok {
			removed = append(removed, p)
		}
	}
	return added, removed
}

// shareVirtualSwitchWithProjects shares the Virtual Switch identified by extID
// with each project UUID in projectIDs, waiting for every share task to reach
// the SUCCEEDED state before returning.
func shareVirtualSwitchWithProjects(ctx context.Context, meta interface{}, extID string, projectIDs []string, timeout time.Duration) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	for _, projectExtID := range projectIDs {
		shareReq := networkingVsReq.ShareVirtualSwitchByIdRequest{
			ExtId: utils.StringPtr(extID),
			Body: &networkingConfig.ProjectReference{
				ProjectExtId: utils.StringPtr(projectExtID),
			},
		}

		aJSON, _ := json.MarshalIndent(shareReq, "", " ")
		log.Printf("[DEBUG] VirtualSwitch share payload: %s", string(aJSON))

		resp, err := conn.VirtualSwitchAPI.ShareVirtualSwitchById(ctx, &shareReq)
		if err != nil {
			return diag.Errorf("error while sharing Virtual Switch (%s) with project (%s): %v", extID, projectExtID, err)
		}

		taskRef := resp.Data.GetValue().(import4.TaskReference)
		if err := waitForVirtualSwitchTask(ctx, meta, taskRef.ExtId, timeout); err != nil {
			return diag.Errorf("error waiting for Virtual Switch (%s) share with project (%s): %s", extID, projectExtID, err)
		}

		aJSON, _ = json.MarshalIndent(taskRef, "", " ")
		log.Printf("[DEBUG] VirtualSwitch share task: %s", string(aJSON))
	}
	return nil
}

// unshareVirtualSwitchFromProjects unshares the Virtual Switch identified by
// extID from each project UUID in projectIDs, waiting for every unshare task
// to reach the SUCCEEDED state before returning.
func unshareVirtualSwitchFromProjects(ctx context.Context, meta interface{}, extID string, projectIDs []string, timeout time.Duration) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	for _, projectExtID := range projectIDs {
		unshareReq := networkingVsReq.UnshareVirtualSwitchByIdRequest{
			ExtId: utils.StringPtr(extID),
			Body: &networkingConfig.ProjectReference{
				ProjectExtId: utils.StringPtr(projectExtID),
			},
		}

		aJSON, _ := json.MarshalIndent(unshareReq, "", " ")
		log.Printf("[DEBUG] VirtualSwitch unshare payload: %s", string(aJSON))

		resp, err := conn.VirtualSwitchAPI.UnshareVirtualSwitchById(ctx, &unshareReq)
		if err != nil {
			return diag.Errorf("error while unsharing Virtual Switch (%s) from project (%s): %v", extID, projectExtID, err)
		}

		taskRef := resp.Data.GetValue().(import4.TaskReference)
		if err := waitForVirtualSwitchTask(ctx, meta, taskRef.ExtId, timeout); err != nil {
			return diag.Errorf("error waiting for Virtual Switch (%s) unshare from project (%s): %s", extID, projectExtID, err)
		}

		aJSON, _ = json.MarshalIndent(taskRef, "", " ")
		log.Printf("[DEBUG] VirtualSwitch unshare task: %s", string(aJSON))
	}
	return nil
}

// waitForVirtualSwitchTask blocks until the prism task identified by taskUUID
// reaches a terminal state, returning a non-nil error if it does not succeed
// within the supplied timeout. Shared by the share/unshare helpers.
func waitForVirtualSwitchTask(ctx context.Context, meta interface{}, taskUUID *string, timeout time.Duration) error {
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: timeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
