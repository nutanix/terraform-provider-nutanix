package ndb

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

type dbID string

const dbIDKey dbID = ""

// this method is used to pass the key-value pair to different modules using context to avoid duplicate code.

// NewContext returns a new Context that carries a provided key value
func NewContext(ctx context.Context, dbID dbID) context.Context {
	return context.WithValue(ctx, dbIDKey, dbID)
}

// FromContext extracts a value from a Context
func FromContext(ctx context.Context) (string, bool) {
	databaseID, ok := ctx.Value(dbIDKey).(dbID)
	return string(databaseID), ok
}

func timeMachineInfoSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		MaxItems:    1,
		ForceNew:    true,
		Optional:    true,
		Description: "sample description for time machine info",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "description of time machine's name",
				},
				"description": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "description of time machine's",
				},
				"slaid": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "description of SLA ID.",
				},
				"sla_details": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"primary_sla": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"sla_id": {
											Type:        schema.TypeString,
											Required:    true,
											Description: "description of SLA ID.",
										},
										"nx_cluster_ids": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
				},
				"autotunelogdrive": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "description of autoTuneLogDrive",
				},
				"schedule": {
					Type:        schema.TypeSet,
					MaxItems:    1,
					Required:    true,
					Description: "description of schedule of time machine",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"snapshottimeofday": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "description of schedule of time machine",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"hours": {
											Type:     schema.TypeInt,
											Required: true,
										},
										"minutes": {
											Type:     schema.TypeInt,
											Required: true,
										},
										"seconds": {
											Type:     schema.TypeInt,
											Required: true,
										},
									},
								},
							},
							"continuousschedule": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "description of schedule of time machine",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"enabled": {
											Type:     schema.TypeBool,
											Required: true,
										},
										"logbackupinterval": {
											Type:     schema.TypeInt,
											Required: true,
										},
										"snapshotsperday": {
											Type:     schema.TypeInt,
											Required: true,
										},
									},
								},
							},
							"weeklyschedule": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "description of schedule of time machine",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"enabled": {
											Type:     schema.TypeBool,
											Required: true,
										},
										"dayofweek": {
											Type:     schema.TypeString,
											Required: true,
										},
									},
								},
							},
							"monthlyschedule": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "description of schedule of time machine",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"enabled": {
											Type:     schema.TypeBool,
											Required: true,
										},
										"dayofmonth": {
											Type:     schema.TypeInt,
											Required: true,
										},
									},
								},
							},
							"quartelyschedule": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "description of schedule of time machine",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"enabled": {
											Type:     schema.TypeBool,
											Required: true,
										},
										"startmonth": {
											Type:     schema.TypeString,
											Required: true,
										},
										"dayofmonth": {
											Type:     schema.TypeInt,
											Required: true,
										},
									},
								},
							},
							"yearlyschedule": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "description of schedule of time machine",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"enabled": {
											Type:     schema.TypeBool,
											Required: true,
										},
										"dayofmonth": {
											Type:     schema.TypeInt,
											Required: true,
										},
										"month": {
											Type:     schema.TypeString,
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				"tags": dataSourceEraDBInstanceTags(),
			},
		},
	}
}

func buildTimeMachineSchedule(set *schema.Set) *era.Schedule {
	d := set.List()
	schedMap := d[0].(map[string]interface{})
	sch := &era.Schedule{}

	if cs, ok := schedMap["snapshottimeofday"]; ok && len(cs.([]interface{})) > 0 {
		conSch := &era.Snapshottimeofday{}

		icmps := (cs.([]interface{}))[0].(map[string]interface{})
		if hours, cok := icmps["hours"]; cok {
			conSch.Hours = hours.(int)
		}

		if mins, tok := icmps["minutes"]; tok {
			conSch.Minutes = mins.(int)
		}
		if secs, tok := icmps["seconds"]; tok {
			conSch.Seconds = secs.(int)
		}

		sch.Snapshottimeofday = conSch
	}

	if cs, ok := schedMap["continuousschedule"]; ok && len(cs.([]interface{})) > 0 {
		conSch := &era.Continuousschedule{}

		icmps := (cs.([]interface{}))[0].(map[string]interface{})
		if enabled, cok := icmps["enabled"]; cok {
			conSch.Enabled = enabled.(bool)
		}

		if mins, tok := icmps["logbackupinterval"]; tok {
			conSch.Logbackupinterval = mins.(int)
		}
		if secs, tok := icmps["snapshotsperday"]; tok {
			conSch.Snapshotsperday = secs.(int)
		}

		sch.Continuousschedule = conSch
	}

	if cs, ok := schedMap["weeklyschedule"]; ok && len(cs.([]interface{})) > 0 {
		conSch := &era.Weeklyschedule{}

		icmps := (cs.([]interface{}))[0].(map[string]interface{})
		if hours, cok := icmps["enabled"]; cok {
			conSch.Enabled = hours.(bool)
		}

		if mins, tok := icmps["dayofweek"]; tok {
			conSch.Dayofweek = mins.(string)
		}

		sch.Weeklyschedule = conSch
	}

	if cs, ok := schedMap["monthlyschedule"]; ok && len(cs.([]interface{})) > 0 {
		conSch := &era.Monthlyschedule{}

		icmps := (cs.([]interface{}))[0].(map[string]interface{})
		if hours, cok := icmps["enabled"]; cok {
			conSch.Enabled = hours.(bool)
		}

		if mins, tok := icmps["dayofmonth"]; tok {
			conSch.Dayofmonth = mins.(int)
		}

		sch.Monthlyschedule = conSch
	}

	if cs, ok := schedMap["quartelyschedule"]; ok && len(cs.([]interface{})) > 0 {
		conSch := &era.Quartelyschedule{}

		icmps := (cs.([]interface{}))[0].(map[string]interface{})
		if hours, cok := icmps["enabled"]; cok {
			conSch.Enabled = hours.(bool)
		}

		if mins, tok := icmps["dayofmonth"]; tok {
			conSch.Dayofmonth = mins.(int)
		}
		if secs, tok := icmps["startmonth"]; tok {
			conSch.Startmonth = secs.(string)
		}

		sch.Quartelyschedule = conSch
	}

	if cs, ok := schedMap["yearlyschedule"]; ok && len(cs.([]interface{})) > 0 {
		conSch := &era.Yearlyschedule{}

		icmps := (cs.([]interface{}))[0].(map[string]interface{})
		if hours, cok := icmps["enabled"]; cok {
			conSch.Enabled = hours.(bool)
		}

		if mins, tok := icmps["dayofmonth"]; tok {
			conSch.Dayofmonth = mins.(int)
		}
		if secs, tok := icmps["month"]; tok {
			conSch.Month = secs.(string)
		}

		sch.Yearlyschedule = conSch
	}

	return sch
}

func buildTimeMachineFromResourceData(set *schema.Set) *era.Timemachineinfo {
	d := set.List()
	tMap := d[0].(map[string]interface{})

	out := &era.Timemachineinfo{}

	if tMap != nil {
		if name, ok := tMap["name"]; ok && len(name.(string)) > 0 {
			out.Name = name.(string)
		}

		if des, ok := tMap["description"]; ok && len(des.(string)) > 0 {
			out.Description = des.(string)
		}

		if slaid, ok := tMap["slaid"]; ok && len(slaid.(string)) > 0 {
			out.Slaid = slaid.(string)
		}

		if schedule, ok := tMap["schedule"]; ok && len(schedule.(*schema.Set).List()) > 0 {
			out.Schedule = *buildTimeMachineSchedule(schedule.(*schema.Set))
		}

		if tags, ok := tMap["tags"]; ok && len(tags.([]interface{})) > 0 {
			out.Tags = expandTags(tags.([]interface{}))
		}

		if autotunelogdrive, ok := tMap["autotunelogdrive"]; ok && autotunelogdrive.(bool) {
			out.Autotunelogdrive = autotunelogdrive.(bool)
		}

		if slaDetails, ok := tMap["sla_details"]; ok && len(slaDetails.([]interface{})) > 0 {
			out.SLADetails = buildSLADetails(slaDetails.([]interface{}))
		}
		return out
	}
	return nil
}

func nodesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Computed:    true,
		Description: "Description of nodes",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"properties": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"vmname": {
					Type:     schema.TypeString,
					Required: true,
				},
				"networkprofileid": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"ip_infos": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ip_type": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"ip_addresses": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"computeprofileid": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"nx_cluster_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"dbserverid": { // When createDbServer is false, we can use this field to set the target db server.
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
			},
		},
	}
}

func buildNodesFromResourceData(d *schema.Set) []*era.Nodes {
	argSet := d.List()
	nodes := []*era.Nodes{}

	for _, arg := range argSet {
		val := arg.(map[string]interface{})
		node := &era.Nodes{}

		if prop, ok := val["properties"]; ok {
			node.Properties = expandNodesProperties(prop.(*schema.Set))
		}
		if vmName, ok := val["vmname"]; ok && len(vmName.(string)) > 0 {
			node.Vmname = utils.StringPtr(vmName.(string))
		}
		if networkProfile, ok := val["networkprofileid"]; ok && len(networkProfile.(string)) > 0 {
			node.Networkprofileid = utils.StringPtr(networkProfile.(string))
		}
		if dbServer, ok := val["dbserverid"]; ok && len(dbServer.(string)) > 0 {
			node.DatabaseServerID = utils.StringPtr(dbServer.(string))
		}
		if nxCls, ok := val["nx_cluster_id"]; ok && len(nxCls.(string)) > 0 {
			node.NxClusterID = utils.StringPtr(nxCls.(string))
		}
		if computeProfile, ok := val["computeprofileid"]; ok && len(computeProfile.(string)) > 0 {
			node.ComputeProfileID = utils.StringPtr(computeProfile.(string))
		}
		if infos, ok := val["ip_infos"]; ok && len(infos.([]interface{})) > 0 {
			node.IPInfos = expandIPInfos(infos.([]interface{}))
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func actionArgumentsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "description of action arguments",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
			},
		},
	}
}

func tryToConvertBool(v interface{}) (bool, bool) {
	str := v.(string)
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false, false
	}
	return b, true
}

func buildActionArgumentsFromResourceData(d *schema.Set, args []*era.Actionarguments) []*era.Actionarguments {
	argSet := d.List()
	for _, arg := range argSet {
		var val interface{}
		val = arg.(map[string]interface{})["value"]
		b, ok := tryToConvertBool(arg.(map[string]interface{})["value"])
		if ok {
			val = b
		}

		args = append(args, &era.Actionarguments{
			Name:  arg.(map[string]interface{})["name"].(string),
			Value: val,
		})
	}
	return args
}

func buildSLADetails(pr []interface{}) *era.SLADetails {
	if len(pr) > 0 {
		res := &era.SLADetails{}

		for _, v := range pr {
			val := v.(map[string]interface{})

			if priSLA, pok := val["primary_sla"]; pok {
				res.PrimarySLA = expandPrimarySLA(priSLA.([]interface{}))
			}
		}
		return res
	}
	return nil
}

func expandPrimarySLA(pr []interface{}) *era.PrimarySLA {
	if len(pr) > 0 {
		out := &era.PrimarySLA{}

		for _, v := range pr {
			val := v.(map[string]interface{})

			if slaid, ok := val["sla_id"]; ok {
				out.SLAID = utils.StringPtr(slaid.(string))
			}

			if nxcls, ok := val["nx_cluster_ids"]; ok {
				res := make([]*string, 0)
				nxclster := nxcls.([]interface{})

				for _, v := range nxclster {
					res = append(res, utils.StringPtr(v.(string)))
				}
				out.NxClusterIds = res
			}
		}
		return out
	}
	return nil
}

func expandNodesProperties(pr *schema.Set) []*era.NodesProperties {
	argSet := pr.List()

	out := make([]*era.NodesProperties, 0)
	for _, arg := range argSet {
		var val interface{}
		val = arg.(map[string]interface{})["value"]
		b, ok := tryToConvertBool(arg.(map[string]interface{})["value"])
		if ok {
			val = b
		}

		out = append(out, &era.NodesProperties{
			Name:  arg.(map[string]interface{})["name"].(string),
			Value: val,
		})
	}
	return out
}

func expandIPInfos(pr []interface{}) []*era.IPInfos {
	if len(pr) > 0 {
		IPInfos := make([]*era.IPInfos, 0)

		for _, v := range pr {
			val := v.(map[string]interface{})
			IPInfo := &era.IPInfos{}

			if ipType, ok := val["ip_type"]; ok {
				IPInfo.IPType = utils.StringPtr(ipType.(string))
			}

			if addr, ok := val["ip_addresses"]; ok {
				res := make([]*string, 0)
				ips := addr.([]interface{})

				for _, v := range ips {
					res = append(res, utils.StringPtr(v.(string)))
				}
				IPInfo.IPAddresses = res
			}

			IPInfos = append(IPInfos, IPInfo)
		}
		return IPInfos
	}
	return nil
}
