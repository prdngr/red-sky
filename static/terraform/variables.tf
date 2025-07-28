variable "aws_profile" {
  type        = string
  default     = "default"
  description = "The AWS profile to use for deployment"
}

variable "aws_region" {
  type        = string
  default     = "eu-central-1"
  description = "The AWS region to deploy in"
}

variable "deployment_id" {
  type        = string
  description = "A unique identifier of the deployment"
}

variable "key_directory" {
  type        = string
  description = "The directory to store the SSH private key in"
}

variable "allowed_ip" {
  type        = string
  default     = null
  description = "(Optional) the allow-listed IP address for ingress traffic"
}

variable "deployment_type" {
  type        = string
  description = "Type of deployment (nessus or kali)"
  validation {
    condition     = contains(["nessus", "kali"], var.deployment_type)
    error_message = "Deployment type must be either 'nessus' or 'kali'"
  }
}
