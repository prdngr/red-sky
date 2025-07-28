output "instance_ip" {
  value       = aws_instance.this.public_ip
  description = "The public IP of the EC2 instance"
}

output "ssh_key_file" {
  value       = local_sensitive_file.this.filename
  description = "The filename of the SSH private key"
}
