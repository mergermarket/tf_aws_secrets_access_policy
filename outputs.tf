output "policy_arn" {
  value = "${aws_iam_policy.secrets_policy.arn}"
}
