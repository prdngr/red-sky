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

module "common" {
  source = "./modules/common"

  deployment_id = var.deployment_id
  key_directory = var.key_directory

  instance_type = "m5.xlarge"
  ami_id        = local.ami_ids[var.deployment_type]
  user_data     = local.user_data[var.deployment_type]
  ingress_rules = local.ingress_rules[var.deployment_type]
}

module "cdn" {
  source = "./modules/cdn"

  count  = var.deployment_type == "c2" ? 1 : 0
  origin_domain = module.common.instance_fqdn
}
