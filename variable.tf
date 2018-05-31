variable "component" {
  type        = "string"
  description = "The component which these secrets belong to"
}

variable "environment" {
  type        = "string"
  description = "The environment which these secrets are for"
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}
