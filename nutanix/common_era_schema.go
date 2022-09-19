package nutanix

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
)

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
					Required:    true,
					Description: "description of time machine's",
				},

				"slaid": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "description of SLA ID.",
				},

				"autotunelogdrive": {
					Type:        schema.TypeBool,
					Required:    true,
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
								Type:        schema.TypeMap,
								Required:    true,
								Description: "description of schedule of time machine",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},

							"continuousschedule": {
								Type:        schema.TypeMap,
								Required:    true,
								Description: "description of schedule of time machine",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},

							"weeklyschedule": {
								Type:        schema.TypeMap,
								Required:    true,
								Description: "description of schedule of time machine",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},

							"monthlyschedule": {
								Type:        schema.TypeMap,
								Required:    true,
								Description: "description of schedule of time machine",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},

							"quartelyschedule": {
								Type:        schema.TypeMap,
								Required:    true,
								Description: "description of schedule of time machine",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},

							"yearlyschedule": {
								Type:        schema.TypeMap,
								Required:    true,
								Description: "description of schedule of time machine",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
				"tags": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					Description: "description of schedule of time machine",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func ConvToInt(s interface{}) int {
	str := s.(string)
	i, _ := strconv.Atoi(str)
	return i
}

func ConvToBool(s interface{}) bool {
	str := s.(string)
	b, _ := strconv.ParseBool(str)
	return b
}

func buildTimeMachineSchedule(set *schema.Set) *era.Schedule {
	d := set.List()
	schedMap := d[0].(map[string]interface{})
	log.Printf("%T", schedMap["snapshottimeofday"].(map[string]interface{})["hours"])
	return &era.Schedule{
		Snapshottimeofday: &era.Snapshottimeofday{
			Hours:   ConvToInt(schedMap["snapshottimeofday"].(map[string]interface{})["hours"]),
			Minutes: ConvToInt(schedMap["snapshottimeofday"].(map[string]interface{})["minutes"]),
			Seconds: ConvToInt(schedMap["snapshottimeofday"].(map[string]interface{})["seconds"]),
		},
		Continuousschedule: &era.Continuousschedule{
			Enabled:           ConvToBool(schedMap["continuousschedule"].(map[string]interface{})["enabled"]),
			Logbackupinterval: ConvToInt(schedMap["continuousschedule"].(map[string]interface{})["logbackupinterval"]),
			Snapshotsperday:   ConvToInt(schedMap["continuousschedule"].(map[string]interface{})["snapshotsperday"]),
		},
		Weeklyschedule: &era.Weeklyschedule{
			Enabled:   ConvToBool(schedMap["weeklyschedule"].(map[string]interface{})["enabled"]),
			Dayofweek: schedMap["weeklyschedule"].(map[string]interface{})["dayofweek"].(string),
		},
		Monthlyschedule: &era.Monthlyschedule{
			Enabled:    ConvToBool(schedMap["monthlyschedule"].(map[string]interface{})["enabled"]),
			Dayofmonth: ConvToInt(schedMap["monthlyschedule"].(map[string]interface{})["dayofmonth"]),
		},
		Quartelyschedule: &era.Quartelyschedule{
			Enabled:    ConvToBool(schedMap["quartelyschedule"].(map[string]interface{})["enabled"]),
			Startmonth: schedMap["quartelyschedule"].(map[string]interface{})["startmonth"].(string),
			Dayofmonth: ConvToInt(schedMap["quartelyschedule"].(map[string]interface{})["dayofmonth"]),
		},
		Yearlyschedule: &era.Yearlyschedule{
			Enabled:    ConvToBool(schedMap["yearlyschedule"].(map[string]interface{})["enabled"]),
			Dayofmonth: ConvToInt(schedMap["yearlyschedule"].(map[string]interface{})["dayofmonth"]),
			Month:      schedMap["yearlyschedule"].(map[string]interface{})["month"].(string),
		},
	}
}

func buildTimeMachineFromResourceData(set *schema.Set) *era.Timemachineinfo {
	d := set.List()
	tMap := d[0].(map[string]interface{})
	return &era.Timemachineinfo{
		Name:             tMap["name"].(string),
		Description:      tMap["description"].(string),
		Slaid:            tMap["slaid"].(string),
		Schedule:         *buildTimeMachineSchedule(tMap["schedule"].(*schema.Set)), // NULL Pointer check
		Tags:             tMap["tags"].(*schema.Set).List(),
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
					Type:        schema.TypeSet,
					Optional:    true,
					ConfigMode:  schema.SchemaConfigModeAttr,
					Description: "",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"vmname": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
					ConfigMode:  schema.SchemaConfigModeAttr,
					Default:     "",
				},
				"networkprofileid": {
					Type:        schema.TypeString,
					Description: "",
					Optional:    true,
					ConfigMode:  schema.SchemaConfigModeAttr,
					Default:     "",
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
