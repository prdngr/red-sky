data "aws_ssm_parameter" "nessus" {
  name = "/aws/service/marketplace/prod-nkomv2twhnmiq/latest"
}

data "aws_ssm_parameter" "kali" {
  name = "/aws/service/marketplace/804fcc46-63fc-4eb6-85a1-50e66d6c7215/latest"
}

data "aws_ec2_managed_prefix_list" "cloudfront" {
  name = "com.amazonaws.global.cloudfront.origin-facing"
}
