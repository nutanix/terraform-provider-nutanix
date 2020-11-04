package nutanix

import (
	"fmt"
	"log"
	"time"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	// UserKind Represents kind of resource
	UserKind = "user"
)

var (
	userDelay      = 10 * time.Second
	userMinTimeout = 3 * time.Second
)

func resourceNutanixUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixUserCreate,
		//Read:   resourceNutanixUserRead,
		//Update: resourceNutanixUserUpdate,
		//Delete: resourceNutanixUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_hash": {
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
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"directory_service_user": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_principal_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							//ValidateFunc: validation.StringInSlice([]string{"role"}, false),
						},
						"directory_service_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:    schema.TypeString,
										Default: "directory_service",
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"identity_provider_user": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							//ValidateFunc: validation.StringInSlice([]string{"role"}, false),
						},
						"identity_provider_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:    schema.TypeString,
										Default: "identity_provider",
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"name": {
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

func resourceNutanixUserCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating User: %s", d.Get("name").(string))
	client := meta.(*Client)
	conn := client.API
	timeout := client.WaitTimeout

	if client.WaitTimeout == 0 {
		timeout = 10
	}

	request := &v3.UserIntentInput{}

	metadata := &v3.Metadata{}

	if err := getMetadataAttributes(d, metadata, "user"); err != nil {
		return err
	}

	spec := &v3.UserSpec{
		Resources: &v3.UserResources{
			DirectoryServiceUser: expandDirectoryServiceUser(d),
			IdentityProviderUser: expandIdentityProviderUser(d),
		},
	}

	request.Metadata = metadata
	request.Spec = spec

	// Make request to the API
	resp, err := conn.V3.CreateUser(request)
	if err != nil {
		return fmt.Errorf("error creating Nutanix User: %+v", err)
	}

	UUID := *resp.Metadata.UUID
	// set terraform state
	d.SetId(UUID)

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    time.Duration(timeout) * time.Minute,
		Delay:      userDelay,
		MinTimeout: userMinTimeout,
	}

	if _, errw := stateConf.WaitForState(); errw != nil {
		// delErr := resourceNutanixUserDelete(d, meta)
		// if delErr != nil {
		// 	return fmt.Errorf("error waiting for image (%s) to delete in creation: %s", d.Id(), delErr)
		// }
		d.SetId("")
		return fmt.Errorf("error waiting for user (%s) to create: %s", UUID, errw)
	}

	return nil
	//return resourceNutanixUserRead(d, meta)
}

// func resourceNutanixUserRead(d *schema.ResourceData, meta interface{}) error {
// 	log.Printf("[DEBUG] Reading Image: %s", d.Get("name").(string))

// 	// Get client connection
// 	conn := meta.(*Client).API
// 	uuid := d.Id()

// 	// Make request to the API
// 	resp, err := conn.V3.GetImage(uuid)
// 	if err != nil {
// 		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
// 			d.SetId("")
// 		}
// 		return fmt.Errorf("error reading image UUID (%s) with error %s", uuid, err)
// 	}

// 	m, c := setRSEntityMetadata(resp.Metadata)

// 	if err = d.Set("metadata", m); err != nil {
// 		return fmt.Errorf("error setting metadata for image UUID(%s), %s", d.Id(), err)
// 	}
// 	if err = d.Set("categories", c); err != nil {
// 		return fmt.Errorf("error setting categories for image UUID(%s), %s", d.Id(), err)
// 	}

// 	if err = d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
// 		return fmt.Errorf("error setting owner_reference for image UUID(%s), %s", d.Id(), err)
// 	}
// 	d.Set("api_version", utils.StringValue(resp.APIVersion))
// 	d.Set("name", utils.StringValue(resp.Status.Name))
// 	d.Set("description", utils.StringValue(resp.Status.Description))

// 	if err = d.Set("availability_zone_reference", flattenReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
// 		return fmt.Errorf("error setting owner_reference for image UUID(%s), %s", d.Id(), err)
// 	}
// 	if err = flattenClusterReference(resp.Status.ClusterReference, d); err != nil {
// 		return fmt.Errorf("error setting cluster_uuid or cluster_name for image UUID(%s), %s", d.Id(), err)
// 	}

// 	if err = d.Set("state", resp.Status.State); err != nil {
// 		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
// 	}

// 	if err = d.Set("image_type", resp.Status.Resources.ImageType); err != nil {
// 		return fmt.Errorf("error setting image_type for image UUID(%s), %s", d.Id(), err)
// 	}

// 	if err = d.Set("source_uri", resp.Status.Resources.SourceURI); err != nil {
// 		return fmt.Errorf("error setting source_uri for image UUID(%s), %s", d.Id(), err)
// 	}

// 	if err = d.Set("size_bytes", resp.Status.Resources.SizeBytes); err != nil {
// 		return fmt.Errorf("error setting size_bytes for image UUID(%s), %s", d.Id(), err)
// 	}

// 	checksum := make(map[string]string)
// 	if resp.Status.Resources.Checksum != nil {
// 		checksum["checksum_algorithm"] = utils.StringValue(resp.Status.Resources.Checksum.ChecksumAlgorithm)
// 		checksum["checksum_value"] = utils.StringValue(resp.Status.Resources.Checksum.ChecksumValue)
// 	}

// 	if err = d.Set("checksum", checksum); err != nil {
// 		return fmt.Errorf("error setting checksum for image UUID(%s), %s", d.Id(), err)
// 	}

// 	version := make(map[string]string)
// 	if resp.Status.Resources.Version != nil {
// 		version["product_version"] = utils.StringValue(resp.Status.Resources.Version.ProductVersion)
// 		version["product_name"] = utils.StringValue(resp.Status.Resources.Version.ProductName)
// 	}

// 	if err = d.Set("version", version); err != nil {
// 		return fmt.Errorf("error setting version for image UUID(%s), %s", d.Id(), err)
// 	}

// 	uriList := make([]string, 0, len(resp.Status.Resources.RetrievalURIList))
// 	for _, uri := range resp.Status.Resources.RetrievalURIList {
// 		uriList = append(uriList, utils.StringValue(uri))
// 	}

// 	if err = d.Set("retrieval_uri_list", uriList); err != nil {
// 		return fmt.Errorf("error setting retrieval_uri_list for image UUID(%s), %s", d.Id(), err)
// 	}

// 	return nil
// }

// func resourceNutanixUserUpdate(d *schema.ResourceData, meta interface{}) error {
// 	client := meta.(*Client)
// 	conn := client.API
// 	timeout := client.WaitTimeout

// 	if client.WaitTimeout == 0 {
// 		timeout = 10
// 	}

// 	// get state
// 	request := &v3.ImageIntentInput{}
// 	metadata := &v3.Metadata{}
// 	spec := &v3.Image{}
// 	res := &v3.ImageResources{}

// 	response, err := conn.V3.GetImage(d.Id())

// 	if err != nil {
// 		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
// 			d.SetId("")
// 		}
// 		return err
// 	}

// 	if response.Metadata != nil {
// 		metadata = response.Metadata
// 	}

// 	if response.Spec != nil {
// 		spec = response.Spec

// 		if response.Spec.Resources != nil {
// 			res = response.Spec.Resources
// 		}
// 	}

// 	if d.HasChange("categories") {
// 		metadata.Categories = expandCategories(d.Get("categories"))
// 	}

// 	if d.HasChange("owner_reference") {
// 		or := d.Get("owner_reference").(map[string]interface{})
// 		metadata.OwnerReference = validateRef(or)
// 	}

// 	if d.HasChange("project_reference") {
// 		pr := d.Get("project_reference").(map[string]interface{})
// 		metadata.ProjectReference = validateRef(pr)
// 	}

// 	if d.HasChange("name") {
// 		spec.Name = utils.StringPtr(d.Get("name").(string))
// 	}
// 	if d.HasChange("description") {
// 		spec.Description = utils.StringPtr(d.Get("description").(string))
// 	}

// 	if d.HasChange("source_uri") || d.HasChange("checksum") {
// 		if err := getImageResource(d, res); err != nil {
// 			return err
// 		}
// 		spec.Resources = res
// 	}

// 	request.Metadata = metadata
// 	request.Spec = spec

// 	resp, errUpdate := conn.V3.UpdateImage(d.Id(), request)

// 	if errUpdate != nil {
// 		return fmt.Errorf("error updating image(%s) %s", d.Id(), errUpdate)
// 	}

// 	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

// 	// Wait for the Image to be available
// 	stateConf := &resource.StateChangeConf{
// 		Pending:    []string{"QUEUED", "RUNNING"},
// 		Target:     []string{"SUCCEEDED"},
// 		Refresh:    taskStateRefreshFunc(conn, taskUUID),
// 		Timeout:    time.Duration(timeout) * time.Minute,
// 		Delay:      userDelay,
// 		MinTimeout: userMinTimeout,
// 	}

// 	if _, err := stateConf.WaitForState(); err != nil {
// 		delErr := resourceNutanixUserDelete(d, meta)
// 		if delErr != nil {
// 			return fmt.Errorf("error waiting for image (%s) to delete in update: %s", d.Id(), delErr)
// 		}
// 		uuid := d.Id()
// 		d.SetId("")
// 		return fmt.Errorf("error waiting for image (%s) to update: %s", uuid, err)
// 	}

// 	return resourceNutanixUserRead(d, meta)
// }

// func resourceNutanixUserDelete(d *schema.ResourceData, meta interface{}) error {
// 	log.Printf("[DEBUG] Deleting Image: %s", d.Get("name").(string))

// 	client := meta.(*Client)
// 	conn := client.API
// 	timeout := client.WaitTimeout

// 	if client.WaitTimeout == 0 {
// 		timeout = 10
// 	}

// 	UUID := d.Id()

// 	resp, err := conn.V3.DeleteImage(UUID)
// 	if err != nil {
// 		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
// 			d.SetId("")
// 		}
// 		return err
// 	}

// 	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

// 	// Wait for the Image to be available
// 	stateConf := &resource.StateChangeConf{
// 		Pending:    []string{"QUEUED", "RUNNING"},
// 		Target:     []string{"SUCCEEDED"},
// 		Refresh:    taskStateRefreshFunc(conn, taskUUID),
// 		Timeout:    time.Duration(timeout) * time.Minute,
// 		Delay:      userDelay,
// 		MinTimeout: userMinTimeout,
// 	}

// 	if _, err := stateConf.WaitForState(); err != nil {
// 		d.SetId("")
// 		return fmt.Errorf("error waiting for image (%s) to delete: %s", d.Id(), err)
// 	}

// 	d.SetId("")
// 	return nil
// }

func expandDirectoryServiceUser(d *schema.ResourceData) *v3.DirectoryServiceUser {
	directoryServiceUserState, ok := d.GetOk("directory_service_user")
	if !ok {
		return nil
	}

	directoryServiceUserMap := directoryServiceUserState.(*schema.Set).List()[0].(map[string]interface{})
	directoryServiceUser := &v3.DirectoryServiceUser{}

	if upn, ok := directoryServiceUserMap["user_principal_name"]; ok {
		directoryServiceUser.UserPrincipalName = utils.StringPtr(upn.(string))
	}

	if dpr, ok := directoryServiceUserMap["directory_service_reference"]; ok {
		directoryServiceUser.DirectoryServiceReference = expandReference(dpr.(*schema.Set).List()[0].(map[string]interface{}))
	}

	return directoryServiceUser
}

func expandIdentityProviderUser(d *schema.ResourceData) *v3.IdentityProvider {
	identityProviderState, ok := d.GetOk("directory_service_user")
	if !ok {
		return nil
	}

	identiryProviderMap := identityProviderState.(*schema.Set).List()[0].(map[string]interface{})
	identiryProvider := &v3.IdentityProvider{}

	if username, ok := identiryProviderMap["username"]; ok {
		identiryProvider.Username = utils.StringPtr(username.(string))
	}

	if ipr, ok := identiryProviderMap["identity_provider_reference"]; ok {
		identiryProvider.IdentityProviderReference = expandReference(ipr.(*schema.Set).List()[0].(map[string]interface{}))
	}

	return identiryProvider
}
