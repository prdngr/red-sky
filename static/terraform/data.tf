data "aws_ami" "nessus" {
  most_recent = true

  filter {
    name   = "name"
    values = ["Nessus *"]
  }

  filter {
    name   = "product-code"
    values = ["8fn69npzmbzcs4blc4583jd0y"]
  }

  owners = ["aws-marketplace"]
}

data "aws_ami" "kali" {
  most_recent = true

  filter {
    name   = "name"
    values = ["kali-last-snapshot-amd64-*"]
  }

  filter {
    name   = "product-code"
    values = ["7lgvy7mt78lgoi4lant0znp5h"]
  }

  owners = ["aws-marketplace"]
}

data "aws_ec2_managed_prefix_list" "cloudfront" {
  name = "com.amazonaws.global.cloudfront.origin-facing"
}
