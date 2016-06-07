package aws

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func defaultInstanceHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["instance_count"].(int)))

	return hashcode.String(buf.String())
}

func resourceAwsElasticMapReduceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElasticMapReduceCreate,
		Read:   resourceAwsElasticMapReduceRead,
		Update: resourceAwsElasticMapReduceUpdate,
		Delete: resourceAwsElasticMapReduceDelete,

		Schema: map[string]*schema.Schema{
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"job_flow_role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"release": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Set:      defaultInstanceHash,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_count": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAwsElasticMapReduceCreate(d *schema.ResourceData, meta interface{}) error {
	emrconn := meta.(*AWSClient).emrconn

	clusterName := d.Get("cluster_name").(string)
	serviceRole := d.Get("service_role").(string)
	jobFlowRole := d.Get("job_flow_role").(string)

	instances := d.Get("instances").(*schema.Set).List()[0].(map[string]interface{})
	instanceCount := instances["instance_count"].(int)

	req := &emr.RunJobFlowInput{
		Name: aws.String(clusterName),
		Instances: &emr.JobFlowInstancesConfig{
			InstanceCount:               aws.Int64(int64(instanceCount)),
			KeepJobFlowAliveWhenNoSteps: aws.Bool(true),
			MasterInstanceType:          aws.String("m1.large"),
			SlaveInstanceType:           aws.String("m1.large"),
			TerminationProtected:        aws.Bool(false),
		},
		ServiceRole:       aws.String(serviceRole),
		JobFlowRole:       aws.String(jobFlowRole),
		VisibleToAllUsers: aws.Bool(true),
	}

	if v, ok := d.GetOk("release"); ok {
		req.ReleaseLabel = aws.String(v.(string))
	}

	resp, err := emrconn.RunJobFlow(req)

	if err != nil {
		return fmt.Errorf("Error creating EMR cluster: %s", err)
	}

	d.SetId(*resp.JobFlowId)
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
