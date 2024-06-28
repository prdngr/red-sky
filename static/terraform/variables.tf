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
  description = "(Optional) the allow-listed IP address for ingress traffic"
}

variable "key_directory" {
  type        = string
  description = "The directory to store the SSH private key in"
}

variable "nessus_username" {
  type        = string
  description = "The Nessus admin username"
}

variable "nessus_password" {
  type        = string
  sensitive   = true
  description = "The Nessus admin password"
}

variable "nessus_activiation_code" {
  type        = string
  sensitive   = true
  default = "activiation-code"
  description = "(Optional) the Nessus activiation code"
}
