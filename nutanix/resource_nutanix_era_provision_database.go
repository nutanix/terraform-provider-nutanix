package nutanix

import (
	"context"
	"encoding/json"
	"log"
	"time"

	client "github.com/terraform-providers/terraform-provider-nutanix/client/era"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: createDatabaseInstance, // TODO: Use CreateContext etc functions
		ReadContext:   readDatabaseInstance,
		UpdateContext: updateDatabaseInstance,
		DeleteContext: deleteDatabaseInstance,
		Schema: map[string]*schema.Schema{
			"database_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true, // TODO: Check whether it is required or not
				Description: "represent id of database instance",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "represent id of database instance",
			},

			"databasetype": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "databast_type: Database type description",
			},

			"softwareprofileid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},

			"softwareprofileversionid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},

			"computeprofileid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},

			"networkprofileid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},
			"dbparameterprofileid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},

			"newdbservertimezone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "UTC",
			},

			"nxclusterid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},

			"sshpublickey": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     "",
			},

			"createdbserver": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     true,
			},

			"dbserverid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "database server ID if createDbserver is false.",
				Default:     "",
			},

			"clustered": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     false,
			},

			"autotunestagingdrive": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "databast_type: Database type description",
				Default:     true,
			},

			"nodecount": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "database_type: Database type description",
				Default:     1,
			},

			"actionarguments": actionArgumentsSchema(),

			"timemachineinfo": timeMachineInfoSchema(),

			"nodes": nodesSchema(),

			"properties": {
				Type:        schema.TypeList,
				Description: "List of all the properties",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},

						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
		},
	}
}

func createDatabaseInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("Creating the request!!!")
	req, err := buildEraRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	c := meta.(*client.Client)
	if c == nil {
		return diag.Errorf("client is nil")
	}

	log.Println("Request:\n")

	b, _ := json.Marshal(req)

	log.Println("Json blob: \n")
	log.Println(string(b))

	log.Println("Sending request to server...............")

	resp, err := c.Service.ProvisionDatabase(req)
	if err != nil {
		return diag.Errorf("error while sending request...........:\n %s\n\n", err.Error())
	}
	d.SetId(resp.Entityid)
	log.Println("Response from server:")
	log.Println(resp)

	log.Println("\n\n\n\n")

	b, _ = json.Marshal(resp)

	log.Println("Json blob: \n")
	log.Println(string(b))

	// TODO: Poll for operation by using operation id we get from response.

	// Get Operation ID from response of ProvisionDatabaseResponse and poll for the operation to get completed.
	opID := resp.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := client.GetOperationRequest{
		OperationID: opID,
	}

	log.Printf("polling for operation with id: %s\n", opID)

	// TODO: change following code to retry timeout mechanism provided by terraform to poll for operation
	for { // Have a timeout too depending upon the time it takes for database instance provision
		log.Println("Waiting for 5 seconds.............")
		time.Sleep(5 * time.Second)
		opRes, err := c.Service.GetOperation(opReq)
		if err != nil {
			return diag.Errorf("error occured while polling for operation with id.. %s: \n%s\n\n", opID, err.Error())
		}
		if opRes.Status == "4" || opRes.Status == "5" {
			if opRes.Status == "4" {
				log.Println("operation with id: %s has failed", opRes.ID)
				return diag.Errorf("operation: %v has failed", opRes)
			} else {
				log.Println("operation with id: %s has completed", opRes.ID)
				log.Println("database instance has successfully created")
			}
			break
		}
	}

	// TODO: Remove all stupid debug statements only have valid debug logs and return response values to schema in computed values.
	return readDatabaseInstance(ctx, d, meta)

}

//func getStr(d *schema.ResourceData, key string) string {
//return d.Get(string).(string)
//}

func buildEraRequest(d *schema.ResourceData) (*client.ProvisionDatabaseRequest, error) {
	return &client.ProvisionDatabaseRequest{
		Databasetype:             d.Get("databasetype").(string),
		Name:                     d.Get("name").(string),
		Databasedescription:      d.Get("description").(string),
		Softwareprofileid:        d.Get("softwareprofileid").(string),
		Softwareprofileversionid: d.Get("softwareprofileversionid").(string),
		Computeprofileid:         d.Get("computeprofileid").(string),
		Networkprofileid:         d.Get("networkprofileid").(string),
		Dbparameterprofileid:     d.Get("dbparameterprofileid").(string),
		Newdbservertimezone:      d.Get("newdbservertimezone").(string),
		DatabaseServerID:         d.Get("dbserverid").(string),
		Timemachineinfo:          *buildTimeMachineFromResourceData(d.Get("timemachineinfo").(*schema.Set)),
		Actionarguments:          buildActionArgumentsFromResourceData(d.Get("actionarguments").(*schema.Set)),
		Createdbserver:           d.Get("createdbserver").(bool),
		Nodecount:                d.Get("nodecount").(int),
		Nxclusterid:              d.Get("nxclusterid").(string),
		Sshpublickey:             d.Get("sshpublickey").(string),
		Clustered:                d.Get("clustered").(bool),
		Nodes:                    buildNodesFromResourceData(d.Get("nodes").(*schema.Set)),
		Autotunestagingdrive:     d.Get("autotunestagingdrive").(bool),
	}, nil
}

func readDatabaseInstance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	if c == nil {
		return diag.Errorf("client is nil")
	}

	databaseInstanceID := d.Id()

	res, err := c.Service.GetDatabaseInstance(databaseInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if res != nil {
		if err = d.Set("description", res.Description); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("name", res.Name); err != nil {
			return diag.FromErr(err)
		}

		props := []interface{}{}
		for _, prop := range res.Properties {
			props = append(props, map[string]interface{}{
				"name":  prop.Name,
				"value": prop.Value,
			})
		}
		if err := d.Set("properties", props); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func updateDatabaseInstance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	if c == nil {
		return diag.Errorf("client is nil")
	}

	dbID := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	updateReq := client.UpdateDatabaseRequest{
		Name:             name,
		Description:      description,
		Tags:             []interface{}{},
		Resetname:        true,
		Resetdescription: true,
		Resettags:        true,
	}

	res, err := c.Service.UpdateDatabase(&updateReq, dbID)
	if err != nil {
		return diag.FromErr(err)
	}

	if res != nil {
		if err = d.Set("description", res.Description); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("name", res.Name); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func deleteDatabaseInstance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	if c == nil {
		return diag.Errorf("client is nil")
	}

	dbID := d.Id()

	req := client.DeleteDatabaseRequest{
		Delete:               false,
		Remove:               true,
		Softremove:           false,
		Forced:               false,
		Deletetimemachine:    false,
		Deletelogicalcluster: true,
	}
	res, err := c.Service.DeleteDatabase(&req, dbID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Operation to delete instance with id %s has started, operation id: %s", dbID, res.Operationid)
	// TODO: Use retry timeout mechanism provided by terraform to poll for operation

	return nil
}
