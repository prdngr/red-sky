variable "aws_region" {
  type        = string
  default     = "eu-central-1"
  description = "The AWS region to deploy in"
}

variable "deployment_name" {
  type        = string
  description = "A unique identifier of the deployment"
}

variable "allowed_ip" {
  type        = string
  default     = null
  description = "If specified, allow-lists the IP address for ingress traffic"
}
