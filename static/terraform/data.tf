data "aws_ssm_parameter" "nessus_latest" {
  name = "/aws/service/marketplace/prod-nkomv2twhnmiq/latest"
}

data "aws_ami" "nessus" {
  owners = ["aws-marketplace"]

  filter {
    name   = "image-id"
    values = [data.aws_ssm_parameter.nessus_latest.insecure_value]
  }
}

data "aws_ami" "kali" {
  most_recent = true
  owners      = ["aws-marketplace"]

  filter {
    name   = "name"
    values = ["*804fcc46-63fc-4eb6-85a1-50e66d6c7215*"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "state"
    values = ["available"]
  }
}

data "aws_ec2_managed_prefix_list" "cloudfront" {
  name = "com.amazonaws.global.cloudfront.origin-facing"
}
