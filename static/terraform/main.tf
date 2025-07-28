provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile

  default_tags {
    tags = {
      Deployment = var.deployment_id
      ManagedBy  = "Terraform"
    }
  }
}

locals {
  ami_ids = {
    nessus = data.aws_ami.nessus.id
    kali   = data.aws_ami.kali.id
  }

  ingress_rules = {
    nessus = [
      {
        port             = 22
        cidr_ipv4        = var.allowed_ip == null ? "0.0.0.0/0" : "${var.allowed_ip}/32"
        description      = "Allow SSH access"
        deployment_types = ["nessus", "kali", "kali-c2"]
      },
      {
        port             = 8834
        cidr_ipv4        = var.allowed_ip == null ? "127.0.0.1/32" : "${var.allowed_ip}/32"
        description      = "Allow Nessus interface access"
        deployment_types = ["nessus", "kali", "kali-c2"]
      }
    ],
    kali = [
      {
        port             = 22
        cidr_ipv4        = var.allowed_ip == null ? "0.0.0.0/0" : "${var.allowed_ip}/32"
        description      = "Allow SSH access"
        deployment_types = ["nessus", "kali", "kali-c2"]
      },
    ]
  }
}

module "common" {
  source = "./modules/common"

  deployment_id = var.deployment_id
  key_directory = var.key_directory

  instance_type = "m5.xlarge"
  ami_id        = local.ami_ids[var.deployment_type]
  ingress_rules = local.ingress_rules[var.deployment_type]
}
