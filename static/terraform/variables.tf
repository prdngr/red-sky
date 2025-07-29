variable "aws_profile" {
  type        = string
  description = "The AWS profile to use for deployment"
}

variable "aws_region" {
  type        = string
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
  description = "The allow-listed IP address for ingress traffic"
}

variable "deployment_type" {
  type        = string
  description = "Type of deployment (nessus, kali, or c2)"
  validation {
    condition     = contains(["nessus", "kali", "c2"], var.deployment_type)
    error_message = "Deployment type must be either 'nessus', 'kali', or 'c2'"
  }
}
