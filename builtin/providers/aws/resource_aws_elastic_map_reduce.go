package aws

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsElasticMapReduce() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElasticMapReduceCreate,
		Read:   resourceAwsElasticMapReduceRead,
		Update: resourceAwsElasticMapReduceUpdate,
		Delete: resourceAwsElasticMapReduceDelete,

		Schema: map[string]*schema.Schema{},
	}
}

func resourceAwsElasticMapReduceCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticMapReduceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticMapReduceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticMapReduceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
