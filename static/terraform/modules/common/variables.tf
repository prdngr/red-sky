variable "deployment_id" {
  type        = string
  description = "A unique identifier of the deployment"
}

variable "key_directory" {
  type        = string
  description = "The directory to store the SSH private key in"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type"
}

variable "ami_id" {
  type        = string
  description = "AMI ID to deploy"
}

variable "user_data" {
  type        = string
  description = "User data script to run on instance launch"
}

variable "ingress_rules" {
  type = list(object({
    port           = number
    cidr_ipv4      = optional(string)
    prefix_list_id = optional(string)
    description    = string
  }))
  description = "List of security group ingress rules"
}
