package main

import (
    "fmt"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

const PLAN_FILE = "plan"

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteDummyProviderConfig() {
	config := `
variable "sts_endpoint" {}

provider "aws" {
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_get_ec2_platforms      = true
  skip_region_validation      = true
  skip_requesting_account_id  = true
  max_retries                 = 1
  access_key                  = "a"
  secret_key                  = "a"
  region                      = "eu-west-1"

  endpoints {
    sts = "${var.sts_endpoint}"
  }
}
`
	err := ioutil.WriteFile("provider.tf", []byte(config), 0644)
	Must(err)
}

func RunTerraformCommand(args ...string) {
    cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), "TF_IN_AUTOMATION=1")
	output, err := cmd.Output()
    if exiterr, ok := err.(*exec.ExitError); ok {
        fmt.Printf("Output from %s: %s\n", args[0], output)
        fmt.Printf("stderr: %s\n", exiterr.Stderr)
    }
    Must(err)
}

func ReadTerraformPlan(planFilePath string) *terraform.Plan {
	fileReader, fileErr := os.Open(planFilePath)
	Must(fileErr)
	plan, planErr := terraform.ReadPlan(fileReader)
	Must(planErr)
	return plan
}

func Setup(tfargs ...string) *TestingPlan {
	WriteDummyProviderConfig()
    RunTerraformCommand("terraform", "init")
    basePlanArgs := []string{"terraform", "plan", "-out", PLAN_FILE}
    tfargs = append(basePlanArgs, tfargs...)
	RunTerraformCommand(
        tfargs...
    )
	plan := ReadTerraformPlan(PLAN_FILE)
	return &TestingPlan{plan}
}

type TestingPlan struct {
	Plan *terraform.Plan
}

func (p *TestingPlan) AssertResource(t *testing.T, resourcePath string) {
	t.Helper()

    _, ok := p.FindResource(resourcePath)
    if ok {
        t.Logf("Found %s\n", resourcePath)
    } else {
        t.Errorf(
            "Expected to find %s in %s",
            resourcePath, p.Plan.Diff.Modules,
        )
    }
}

func (p *TestingPlan) AssertResourceAttribute(
	t *testing.T,
	resourcePath string,
	attributeName string,
	attributeValue string,
) {
	t.Helper()

    resource, ok := p.FindResource(resourcePath)
    if !ok {
        t.Errorf("Could not find %s in %s", resourcePath, p.Plan.Diff.Modules)
    }
    attr, attrOk := p.FindResourceAttribute(resource, attributeName)
    if !attrOk {
        t.Errorf("Did not find %s in %s\n", attributeName, resourcePath)
    }
    if attr.New == attributeValue {
        t.Logf("Found %s = %s, on %s\n", attributeName, attributeValue, resourcePath)
        return
    } else {
        t.Errorf(
            "Expected %s, got %s for %s attribute on %s resource",
            attributeValue, attr.New, attributeName, attributeValue,
        )
    }
}

func (p *TestingPlan) FindResource(resourcePath string) (
    *terraform.InstanceDiff, bool,
) {
    modules := p.Plan.Diff.Modules
	for _, module := range modules {
		for key, instanceDiff := range module.Resources {
			if key == resourcePath {
				return instanceDiff, true
			}
		}
	}
    return nil, false
}

func (p *TestingPlan) FindResourceAttribute(
    resource *terraform.InstanceDiff, attributeName string,
) (*terraform.ResourceAttrDiff, bool) {
    for attrKey, attr := range resource.Attributes {
        if attrKey == attributeName {
            return attr, true
        }
    }
    return nil, false
}

func TestPolicy(t *testing.T) {
    args := []string{
        "-var", "component=mycomponent",
        "-var", "environment=test",
    }
	plan := Setup(args...)

	plan.AssertResource(t, "aws_iam_policy.secrets_policy")
    plan.AssertResourceAttribute(
        t, "aws_iam_policy.secrets_policy", "name",
        "test-mycomponent-secrets",
    )
    expectedPolicy := `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "secretsmanager:GetSecretValue",
      "Resource": "arn:aws:secretsmanager:eu-west-1:123456789012:secret:mycomponent/test/*"
    }
  ]
}`
    policyResource, _ := plan.FindResource("aws_iam_policy.secrets_policy")
    policy, _ := plan.FindResourceAttribute(policyResource, "policy")

    if policy.New != expectedPolicy {
        t.Errorf("Expected %s, got %s", expectedPolicy, policy.New)
    }
}

func TestPolicyIncludingTeam(t *testing.T) {
    args := []string{
        "-var", "component=mycomponent",
        "-var", "environment=test",
        "-var", "team=someteam",
    }
	plan := Setup(args...)

	plan.AssertResource(t, "aws_iam_policy.secrets_policy")
    plan.AssertResourceAttribute(
        t, "aws_iam_policy.secrets_policy", "name",
        "test-mycomponent-secrets",
    )
    expectedPolicy := `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "secretsmanager:GetSecretValue",
      "Resource": "arn:aws:secretsmanager:eu-west-1:123456789012:secret:someteam/mycomponent/test/*"
    }
  ]
}`
    policyResource, _ := plan.FindResource("aws_iam_policy.secrets_policy")
    policy, _ := plan.FindResourceAttribute(policyResource, "policy")

    if policy.New != expectedPolicy {
        t.Errorf("Expected %s, got %s", expectedPolicy, policy.New)
    }
}
