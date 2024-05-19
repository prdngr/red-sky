variable "aws_region" {
  type        = string
  default     = "eu-central-1"
  description = "The AWS region to deploy in"
}

variable "deployment_name" {
  type        = string
  description = "A unique identifier of the deployment"
}

variable "source_ip" {
  type        = string
  default     = "127.0.0.1"
  description = "The public IP address that will be used to access the deployment"
}

variable "allow_source_ip" {
  type        = bool
  default     = false
  description = "Whether to allow incoming HTTPS traffic from the source IP"
}
