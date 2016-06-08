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
				ForceNew: true,
			},
			"emr_role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ec2_instance_profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"release": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"applications": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
						"termination_protection": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auto_terminate": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
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
	emrRole := d.Get("emr_role").(string)
	ec2InstanceProfile := d.Get("ec2_instance_profile").(string)

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
		ServiceRole:       aws.String(emrRole),
		JobFlowRole:       aws.String(ec2InstanceProfile),
		VisibleToAllUsers: aws.Bool(true),
	}

	if v, ok := d.GetOk("release"); ok {
		req.ReleaseLabel = aws.String(v.(string))
	}

	if v, ok := d.GetOk("auto_terminate"); ok {
		req.Instances.KeepJobFlowAliveWhenNoSteps = aws.Bool(!v.(bool))
	}

	if v, ok := d.GetOk("termination_protection"); ok {
		req.Instances.TerminationProtected = aws.Bool(v.(bool))
	}

	applications := d.Get("applications").(*schema.Set).List()
	if len(applications) > 0 {
		req.Applications = expandApplications(applications)
	}

	resp, err := emrconn.RunJobFlow(req)

	if err != nil {
		return fmt.Errorf("Error creating EMR cluster: %s", err)
	}

	d.SetId(*resp.JobFlowId)
	return nil
}

func resourceAwsElasticMapReduceRead(d *schema.ResourceData, meta interface{}) error {
	emrconn := meta.(*AWSClient).emrconn

	req := &emr.DescribeClusterInput{
		ClusterId: aws.String(d.Id()),
	}

	resp, err := emrconn.DescribeCluster(req)
	if err != nil {
		return fmt.Errorf("Error reading EMR cluster: %s", err)
	}
	fmt.Println(resp)

	cluster := resp.Cluster

	d.Set("cluster_name", cluster.Name)
	d.Set("release", cluster.ReleaseLabel)

	return nil
}

func resourceAwsElasticMapReduceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticMapReduceDelete(d *schema.ResourceData, meta interface{}) error {
	emrconn := meta.(*AWSClient).emrconn

	req := &emr.TerminateJobFlowsInput{
		JobFlowIds: []*string{
			aws.String(d.Id()),
		},
	}

	_, err := emrconn.TerminateJobFlows(req)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func expandApplications(apps []interface{}) []*emr.Application {
	appOut := make([]*emr.Application, 0, len(apps))
	for _, appName := range expandStringList(apps) {
		app := &emr.Application{
			Name: appName,
		}
		appOut = append(appOut, app)
	}
	return appOut
}
