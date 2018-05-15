package nutanix

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func getMetadataAttributes(d *schema.ResourceData, metadata *v3.Metadata, kind string) error {
	m, mok := d.GetOk("metadata")
	metad := m.(map[string]interface{})

	if !mok {
		return fmt.Errorf("please provide metadata required attributes")
	}

	metadata.Kind = utils.String(kind)

	if v, ok := metad["uuid"]; ok && v != "" {
		metadata.UUID = utils.String(v.(string))
	}
	if v, ok := metad["spec_version"]; ok && v != 0 {
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			return err
		}
		metadata.SpecVersion = utils.Int64(int64(i))
	}
	if v, ok := metad["spec_hash"]; ok && v != "" {
		metadata.SpecHash = utils.String(v.(string))
	}
	if v, ok := metad["name"]; ok {
		metadata.Name = utils.String(v.(string))
	}
	if v, ok := d.GetOk("categories"); ok {
		c := v.(map[string]interface{})
		labels := map[string]string{}

		for k, v := range c {
			labels[k] = v.(string)
		}
		metadata.Categories = labels
	}
	if p, ok := d.GetOk("project_reference"); ok {
		pr := p.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.ProjectReference = r
	}
	if o, ok := metad["owner_reference"]; ok {
		or := o.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(or["kind"].(string)),
			UUID: utils.String(or["uuid"].(string)),
		}
		if v1, ok1 := or["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.OwnerReference = r
	}

	return nil
}
