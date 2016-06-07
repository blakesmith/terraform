package aws

import (
	"fmt"
	//	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccAWSEMrCluster_basic(t *testing.T) {
	var jobFlow emr.RunJobFlowOutput
	resource.Test(t, resource.TestCase{
		PreCheck:     nil,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSEmrDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSEmrClusterConfig,
				Check:  testAccCheckAWSEmrClusterExists("aws_elastic_map_reduce_cluster.tf-test-cluster", &jobFlow),
			},
		},
	})
}

func testAccCheckAWSEmrDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckAWSEmrClusterExists(n string, v *emr.RunJobFlowOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		fmt.Printf("Cluster primary id is: %s", rs.Primary.ID)
		return nil
	}
}

var testAccAWSEmrClusterConfig = fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}

resource "aws_iam_role" "service_role" {
  name = "tf-emr-service-role-%s"
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_role" "job_flow_role" {
  name = "tf-emr-job-flow-role-%s"
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_policy_attachment" "service_attach" {
  name = "tf-service-role-attach-%s"
  roles = ["${aws_iam_role.service_role.name}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceRole"
}

resource "aws_iam_policy_attachment" "job_flow_attach" {
  name = "tf-job-flow-role-attach-%s"
  roles = ["${aws_iam_role.job_flow_role.name}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceforEC2Role"
}
  
resource "aws_elastic_map_reduce_cluster" "tf-test-cluster" {
  cluster_name = "tf-emr-%s"
  release = "emr-4.7.0"
  instances {
    instance_count = 2
  }
  service_role = "${aws_iam_role.service_role.name}"
  job_flow_role = "${aws_iam_role.job_flow_role.name}"
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
