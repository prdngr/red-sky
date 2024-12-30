output deployment_id {
  value       = var.deployment_id
  description = "The unique identifier of the deployment"
}

output "instance_id" {
  value       = aws_instance.nessus.id
  description = "The ID of the EC2 instance"
}

output "instance_ip" {
  value       = aws_instance.nessus.public_ip
  description = "The public IP of the EC2 instance"
}

output "ssh_key_file" {
  value       = local_sensitive_file.this.filename
  description = "The filename of the SSH private key"
}
