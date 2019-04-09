resource "aws_iam_policy" "secrets_policy" {
  name        = "${var.environment}-${var.component}-secrets"
  description = "Secrets access"
  policy      = "${data.aws_iam_policy_document.secrets_access.json}"
}

locals {
  arn_base         = "arn:aws:secretsmanager:${data.aws_region.current.name}"
  secret_namespace = "${var.team == "" ? "" : "${var.team}/"}${var.component}/${var.environment}/*"
}

data "aws_iam_policy_document" "secrets_access" {
  statement {
    actions = [
      "secretsmanager:GetSecretValue",
    ]

    resources = [
      "${local.arn_base}:${data.aws_caller_identity.current.account_id}:secret:${local.secret_namespace}",
    ]
  }
}
