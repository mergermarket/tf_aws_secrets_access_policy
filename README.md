# AWS Secrets Manager Access Policy

Grants access to a namespace in Secrets Manager for a component in an environment.

## Example

```hcl
resource "aws_iam_role_policy_attachment" "secrets_attach" {                    
    role       = "${var.iam_role_name}"               
    policy_arn = "${module.secrets_policy.policy_arn}"                          
}                                                                               
                                                                                
module "secrets_policy" {                                                       
    source = "github.com/mergermarket/tf_aws_secrets_access_policy"             
                                                                                
    component   = "my_service"                                 
    environment = "live"                                                  
}
```
