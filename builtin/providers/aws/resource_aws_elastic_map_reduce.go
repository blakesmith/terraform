package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/hashicorp/terraform/helper/schema"
)

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

			// "software": &schema.Schema{
			// 	Type: schema.TypeSet,
			// 	Optional: true,
			// Set: cacheBehavior
			// },
			// "release": &schema.Schema{
			// 	Type: schema.TypeString,
			// 	Optional: true,
			// },
		},
	}
}

func resourceAwsElasticMapReduceCreate(d *schema.ResourceData, meta interface{}) error {
	emrconn := meta.(*AWSClient).emrconn

	params := &emr.RunJobFlowInput{
		Name:       aws.String("MyCluster"), // Required
		AmiVersion: aws.String("3.8"),
		Steps: []*emr.StepConfig{
			{ // Required
			// HadoopJarStep: &emr.HadoopJarStepConfig{ // Required
			// 	Jar: aws.String("XmlString"), // Required
			// 	Args: []*string{
			// 		aws.String("XmlString"), // Required
			// 		// More values...
			// 	},
			// 	MainClass: aws.String("XmlString"),
			// 	Properties: []*emr.KeyValue{
			// 		{ // Required
			// 			Key:   aws.String("XmlString"),
			// 			Value: aws.String("XmlString"),
			// 		},
			// 		// More values...
			// 	},
			// },
			// Name:            aws.String("XmlStringMaxLen256"), // Required
			// ActionOnFailure: aws.String("ActionOnFailure"),
			},
		},
		Instances: &emr.JobFlowInstancesConfig{ // Required
			InstanceCount:               aws.Int64(1),
			KeepJobFlowAliveWhenNoSteps: aws.Bool(true),
			MasterInstanceType:          aws.String("m1.large"),
			SlaveInstanceType:           aws.String("m1.large"),
			TerminationProtected:        aws.Bool(false),
		},
		ServiceRole:       aws.String("EMR_DefaultRole"),
		JobFlowRole:       aws.String("EMR_EC2_DefaultRole"),
		VisibleToAllUsers: aws.Bool(true),
	}
	resp, err := emrconn.RunJobFlow(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}

	// Pretty-print the response data.
	fmt.Println(resp)

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
