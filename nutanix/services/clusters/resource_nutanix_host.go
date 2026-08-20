package clusters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
)

var (
	hostDelay      = 10 * time.Second
	hostMinTimeout = 3 * time.Second
)

func ResourceNutanixHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixHostCreate,
		ReadContext:   resourceNutanixHostRead,
		UpdateContext: resourceNutanixHostUpdate,
		DeleteContext: resourceNutanixHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Update: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Delete: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the existing Nutanix host to manage. Hosts cannot be created or deleted via Terraform; this resource adopts an existing host to manage its categories.",
			},
			"categories": categoriesSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_driver_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failover_cluster": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipmi": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cpu_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_nics_id_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_cpu_sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"windows_domain": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gpu_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"num_virtual_display_heads": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"assignable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"license_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"num_vgpus_allocated": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pci_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frame_buffer_size_mib": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"numa_node": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_resolution": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fraction": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"guest_driver_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_capacity_hz": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_capacity_mib": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host_disks_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"monitoring_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_cpu_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rackable_unit_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"controller_vm": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"block": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNutanixHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	hostID := d.Get("host_id").(string)

	host, err := conn.V3.GetHost(hostID)
	if err != nil {
		return diag.Errorf("error retrieving host with ID (%s): %s", hostID, err)
	}

	if v, ok := d.GetOk("categories"); ok {
		if err := updateHostCategories(ctx, d, conn, host, expandCategories(v)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(hostID)

	return resourceNutanixHostRead(ctx, d, meta)
}

func resourceNutanixHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	hostID := d.Id()

	host, err := conn.V3.GetHost(hostID)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error retrieving host with ID (%s): %s", hostID, err)
	}

	if err := d.Set("host_id", hostID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", host.Status.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpu_driver_version", host.Status.Resources.GPUDriverVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("failover_cluster", flattenFailOverCluster(host.Status.Resources.FailoverCluster)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ipmi", flattenIMPI(host.Status.Resources.IPMI)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_model", host.Status.Resources.CPUModel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_nics_id_list", host.Status.Resources.HostNicsIDList); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_cpu_sockets", host.Status.Resources.NumCPUSockets); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("windows_domain", flattenWindowsDomain(host.Status.Resources.WindowsDomain)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpu_list", flattenGpuList(host.Status.Resources.GPUList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("serial_number", host.Status.Resources.SerialNumber); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_capacity_hz", host.Status.Resources.CPUCapacityHZ); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_capacity_mib", host.Status.Resources.MemoryVapacityMib); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_disks_reference_list", flattenReferenceList(host.Status.Resources.HostDisksReferenceList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("monitoring_state", host.Status.Resources.MonitoringState); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor", flattenHypervisor(host.Status.Resources.Hypervisor)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_type", host.Status.Resources.HostType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_cpu_cores", host.Status.Resources.NumCPUCores); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rackable_unit_reference", flattenReference(host.Status.Resources.RackableUnitReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_vm", flattenControllerVM(host.Status.Resources.ControllerVM)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("block", flattenBlock(host.Status.Resources.Block)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", flattenReference(host.Status.ClusterReference)); err != nil {
		return diag.FromErr(err)
	}

	m, c := setRSEntityMetadata(host.Metadata)

	if err := d.Set("api_version", host.APIVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", c); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_reference", flattenReferenceValues(host.Metadata.ProjectReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_reference", flattenReferenceValues(host.Metadata.OwnerReference)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	hostID := d.Id()

	if d.HasChange("categories") {
		host, err := conn.V3.GetHost(hostID)
		if err != nil {
			return diag.Errorf("error retrieving host with ID (%s): %s", hostID, err)
		}

		if err := updateHostCategories(ctx, d, conn, host, expandCategories(d.Get("categories"))); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNutanixHostRead(ctx, d, meta)
}

func resourceNutanixHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	hostID := d.Id()

	host, err := conn.V3.GetHost(hostID)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error retrieving host with ID (%s): %s", hostID, err)
	}

	if err := updateHostCategories(ctx, d, conn, host, map[string]string{}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func updateHostCategories(ctx context.Context, d *schema.ResourceData, conn *v3.Client, host *v3.HostResponse, categories map[string]string) error {
	metadata := &v3.Metadata{}
	if host.Metadata != nil {
		metadata = host.Metadata
	}
	metadata.Categories = categories

	request := &v3.HostIntentInput{
		APIVersion: host.APIVersion,
		Metadata:   metadata,
		Spec:       host.Spec,
	}

	resp, err := conn.V3.UpdateHost(*host.Metadata.UUID, request)
	if err != nil {
		return fmt.Errorf("error updating categories for host (%s): %w", *host.Metadata.UUID, err)
	}

	if resp.Status == nil || resp.Status.ExecutionContext == nil || resp.Status.ExecutionContext.TaskUUID == nil {
		return nil
	}

	taskUUID := cast.ToString(resp.Status.ExecutionContext.TaskUUID)
	if taskUUID == "" {
		return nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      hostDelay,
		MinTimeout: hostMinTimeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for host (%s) categories update to finish: %w", *host.Metadata.UUID, err)
	}

	return nil
}
