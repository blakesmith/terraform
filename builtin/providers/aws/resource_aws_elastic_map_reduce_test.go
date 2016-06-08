package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
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
		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster id set")
		}
		conn := testAccProvider.Meta().(*AWSClient).emrconn
		_, err := conn.DescribeCluster(&emr.DescribeClusterInput{
			ClusterId: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return fmt.Errorf("EMR error: %v", err)
		}
		return nil
	}
}

var testAccAWSEmrClusterConfig = fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}

resource "aws_elastic_map_reduce_cluster" "tf-test-cluster" {
  cluster_name = "tf-emr-%s"
  release = "emr-4.7.0"
  applications = ["hive", "hadoop", "pig", "spark", "hue"]
  instances {
    instance_count = 2
    auto_terminate = true
    termination_protection = true
  }
  emr_role = "EMR_DefaultRole"
  ec2_instance_profile = "EMR_EC2_DefaultRole"
}
`, acctest.RandString(10))
