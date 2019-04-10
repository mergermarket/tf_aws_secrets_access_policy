variable "component" {
  type        = "string"
  description = "The component which these secrets belong to"
}

variable "environment" {
  type        = "string"
  description = "The environment which these secrets are for"
}

variable "team" {
  type        = "string"
  description = "The team who owns this secret"
  default     = ""
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}
