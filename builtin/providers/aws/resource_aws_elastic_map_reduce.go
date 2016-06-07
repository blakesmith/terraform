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

	instances := d.Get("instances").(*schema.Set).List()[0].(map[string]interface{})
	instanceCount := instances["instance_count"].(int)

	req := &emr.RunJobFlowInput{
		Name: aws.String(clusterName),
		// Steps: []*emr.StepConfig{
		// 	{ // Required
		// 	// HadoopJarStep: &emr.HadoopJarStepConfig{ // Required
		// 	// 	Jar: aws.String("XmlString"), // Required
		// 	// 	Args: []*string{
		// 	// 		aws.String("XmlString"), // Required
		// 	// 		// More values...
		// 	// 	},
		// 	// 	MainClass: aws.String("XmlString"),
		// 	// 	Properties: []*emr.KeyValue{
		// 	// 		{ // Required
		// 	// 			Key:   aws.String("XmlString"),
		// 	// 			Value: aws.String("XmlString"),
		// 	// 		},
		// 	// 		// More values...
		// 	// 	},
		// 	// },
		// 	// Name:            aws.String("XmlStringMaxLen256"), // Required
		// 	// ActionOnFailure: aws.String("ActionOnFailure"),
		// 	},
		// },
		Instances: &emr.JobFlowInstancesConfig{ // Required
			InstanceCount:               aws.Int64(int64(instanceCount)),
			KeepJobFlowAliveWhenNoSteps: aws.Bool(true),
			MasterInstanceType:          aws.String("m1.large"),
			SlaveInstanceType:           aws.String("m1.large"),
			TerminationProtected:        aws.Bool(false),
		},
		ServiceRole:       aws.String("EMR_DefaultRole"),
		JobFlowRole:       aws.String("EMR_EC2_DefaultRole"),
		VisibleToAllUsers: aws.Bool(true),
	}

	if v, ok := d.GetOk("release"); ok {
		req.ReleaseLabel = aws.String(v.(string))
	}

	resp, err := emrconn.RunJobFlow(req)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}

	// Pretty-print the response data.
	fmt.Println(resp)
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
