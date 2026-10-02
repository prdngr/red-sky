locals {
  ami_ids = {
    nessus = data.aws_ami.nessus.id
    kali   = data.aws_ami.kali.id
    c2     = data.aws_ami.kali.id
  }

  user_data = {
    nessus = null
    kali   = null
    c2     = file("${path.module}/scripts/c2-user-data.sh")
  }

  instance_type = {
    nessus = "m5.xlarge"
    kali   = "m5.large"
    c2     = "m5.large"
  }

  ingress_rules = {
    nessus = [
      {
        port        = 22
        cidr_ipv4   = var.admin_cidr == null ? "0.0.0.0/0" : var.admin_cidr
        description = "Allow SSH access"
      },
      {
        port        = 8834
        cidr_ipv4   = var.admin_cidr == null ? "127.0.0.1/32" : var.admin_cidr
        description = "Allow Nessus interface access"
      }
    ],
    kali = [
      {
        port        = 22
        cidr_ipv4   = var.admin_cidr == null ? "0.0.0.0/0" : var.admin_cidr
        description = "Allow SSH access"
      }
    ],
    c2 = [
      {
        port        = 22
        cidr_ipv4   = var.admin_cidr == null ? "0.0.0.0/0" : var.admin_cidr
        description = "Allow SSH access"
      },
      {
        port        = 7443
        cidr_ipv4   = var.admin_cidr == null ? "127.0.0.1/32" : var.admin_cidr
        description = "Allow Mythic interface access"
      },
      {
        port           = 80
        prefix_list_id = data.aws_ec2_managed_prefix_list.cloudfront.id
        description    = "Allow HTTP access from CloudFront"
      }
    ]
  }
}
