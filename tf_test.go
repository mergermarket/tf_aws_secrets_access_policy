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
    if err != nil {
        fmt.Println(output)
        panic(err)
    }
}

func ReadTerraformPlan(planFilePath string) *terraform.Plan {
	fileReader, fileErr := os.Open(planFilePath)
	Must(fileErr)
	plan, planErr := terraform.ReadPlan(fileReader)
	Must(planErr)
	return plan
}

func Setup() *TestingPlan {
	WriteDummyProviderConfig()
    RunTerraformCommand("terraform", "init")
	RunTerraformCommand(
        "terraform", "plan",
        "-out", PLAN_FILE,
        "-var", "component=a",
        "-var", "environment=b",
    )
	plan := ReadTerraformPlan(PLAN_FILE)
	return &TestingPlan{plan}
}

type TestingPlan struct {
	Plan *terraform.Plan
}

func (p *TestingPlan) AssertResource(t *testing.T, resourcePath string) {
	t.Helper()

    _, ok := FindResource(p.Plan.Diff.Modules, resourcePath)
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

    resource, ok := FindResource(p.Plan.Diff.Modules, resourcePath)
    if !ok {
        t.Errorf("Could not find %s in %s", resourcePath, p.Plan.Diff.Modules)
    }
    for attrKey, attr := range resource.Attributes {
        if attrKey == attributeName {
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
    }
    t.Errorf("Did not find %s in %s\n", attributeName, resourcePath)
}

func FindResource(modules []*terraform.ModuleDiff, resourcePath string) (
    *terraform.InstanceDiff, bool,
) {
	for _, module := range modules {
		for key, instanceDiff := range module.Resources {
			if key == resourcePath {
				return instanceDiff, true
			}
		}
	}
    return nil, false
}

func TestPolicy(t *testing.T) {
	plan := Setup()

	plan.AssertResource(t, "aws_iam_policy.secrets_policy")
    plan.AssertResourceAttribute(t, "aws_iam_policy.secrets_policy", "name", "b-a-secrets")
}
