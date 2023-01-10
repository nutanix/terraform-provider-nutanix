package nutanix

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
)

type dbID string

const dbIDKey dbID = ""

func NewContext(ctx context.Context, dbID dbID) context.Context {
	return context.WithValue(ctx, dbIDKey, dbID)
}

func FromContext(ctx context.Context) (dbID, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the net.IP type assertion returns ok=false for nil.
	dbID, ok := ctx.Value(dbIDKey).(dbID)
	return dbID, ok
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
					Required:    true,
					Description: "description of SLA ID.",
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
	return &era.Timemachineinfo{
		Name:             tMap["name"].(string),
		Description:      tMap["description"].(string),
		Slaid:            tMap["slaid"].(string),
		Schedule:         *buildTimeMachineSchedule(tMap["schedule"].(*schema.Set)), // NULL Pointer check
		Tags:             expandTags(tMap["tags"].([]interface{})),
		Autotunelogdrive: tMap["autotunelogdrive"].(bool),
	}
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
					Type:       schema.TypeSet,
					Optional:   true,
					ConfigMode: schema.SchemaConfigModeAttr,
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
					Type:       schema.TypeString,
					Required:   true,
					ConfigMode: schema.SchemaConfigModeAttr,
				},
				"networkprofileid": {
					Type:       schema.TypeString,
					Required:   true,
					ConfigMode: schema.SchemaConfigModeAttr,
				},
				"dbserverid": { // When createDbServer is false, we can use this field to set the target db server.
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
					ConfigMode:  schema.SchemaConfigModeAttr,
					Default:     "",
				},
			},
		},
	}
}

func buildNodesFromResourceData(d *schema.Set) []*era.Nodes {
	argSet := d.List()
	args := []*era.Nodes{}

	for _, arg := range argSet {
		args = append(args, &era.Nodes{
			Properties:       arg.(map[string]interface{})["properties"].(*schema.Set).List(),
			Vmname:           arg.(map[string]interface{})["vmname"].(string),
			Networkprofileid: arg.(map[string]interface{})["networkprofileid"].(string),
			DatabaseServerID: arg.(map[string]interface{})["dbserverid"].(string),
		})
	}
	return args
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
