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
				Check:  testAccCheckAWSEmrClusterExists("aws_elastic_map_reduce_clust.tf-test-cluster", &jobFlow),
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
resource "aws_security_group" "bar" {
  name = "tf-test-security-group-%03d"
  description = "tf-test-security-group-descr"
  ingress {
    from_port = -1
    to_port = -1
    protocol = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_elastic_map_reduce_cluster" "tf-test-cluster" {
  cluster_name = "tf-emr-%s"
  release = "emr-4.7.0"
  instances {
    instance_count = 2
  }
}
`, acctest.RandInt(), acctest.RandString(10))
