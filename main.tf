resource "aws_iam_policy" "secrets_policy" {                                    
    name        = "${var.environment}-${var.component}-secrets"      
    description = "Secrets access"                                              
    policy      = "${data.aws_iam_policy_document.secrets_access.json}"         
}

data "aws_iam_policy_document" "secrets_access" {                               
  statement {                                                                   
    actions = [                                                                 
      "secretsmanager:GetSecretValue",                                          
    ]                                                                           
                                                                                
    resources = [                                                               
      "arn:aws:secretsmanager:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:secret:${var.component}/${var.environment}/*",
    ]                                                                           
  }                                                                             
}
