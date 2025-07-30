output "deployment_id" {
  value       = var.deployment_id
  description = "The unique identifier of the deployment"
}

output "instance_ip" {
  value       = module.common.instance_ip
  description = "The public IP of the EC2 instance"
}

output "ssh_key_file" {
  value       = module.common.ssh_key_file
  description = "The filename of the SSH private key"
}

output "cloudfront_url" {
  value       = var.deployment_type == "c2" ? module.cdn[0].cloudfront_url : null
  description = "CloudFront distribution URL"
}
