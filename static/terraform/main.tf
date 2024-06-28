provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Workload   = "Nessus on Demand"
      Deployment = var.deployment_name
      ManagedBy  = "Terraform"
    }
  }
}

# ------------------------------------------------------------------------------
# EC2 INSTANCE
# ------------------------------------------------------------------------------

resource "aws_instance" "nessus" {
  ami                    = data.aws_ami.nessus.id
  instance_type          = "m5.xlarge"
  vpc_security_group_ids = [aws_security_group.this.id]
  key_name               = aws_key_pair.this.key_name

  user_data = <<-EOF
    #cloud-config

    write_files:
      - path: /opt/nessus/var/nessus/config.json
        content: |
          {
            "user": {
              "username": "${var.nessus_username}",
              "password": "${var.nessus_password}",
              "role": "system_administrator",
              "type": "local"
            }
          }

    runcmd:
      - while [ $(/opt/nessus/sbin/nessuscli lsuser | wc -l) -ne 1 ]; do sleep 5; done
      - service nessusd stop
      - /opt/nessus/sbin/nessuscli fetch --register ${var.nessus_activiation_code}
      - service nessusd start
  EOF
}

# ------------------------------------------------------------------------------
# SSH KEY PAIR
# ------------------------------------------------------------------------------

resource "tls_private_key" "this" {
  algorithm = "ED25519"
}

resource "aws_key_pair" "this" {
  public_key = tls_private_key.this.public_key_openssh
}

resource "local_sensitive_file" "this" {
  filename        = pathexpand("${var.key_directory}/${var.deployment_name}.pem")
  content         = tls_private_key.this.private_key_openssh
  file_permission = "400"
}

# ------------------------------------------------------------------------------
# SECURITY GROUP
# ------------------------------------------------------------------------------

resource "aws_security_group" "this" {
  description = "Nessus on Demand"
}

# Always allow incoming SSH and outgoing ANY.

resource "aws_vpc_security_group_ingress_rule" "ingress_ssh" {
  security_group_id = aws_security_group.this.id
  from_port         = 22
  to_port           = 22
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_egress_rule" "egress_any" {
  security_group_id = aws_security_group.this.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}

# (Optional) allow-list incoming HTTPS from specified IP.

resource "aws_vpc_security_group_ingress_rule" "ingress_https" {
  count             = var.allowed_ip == null ? 0 : 1
  security_group_id = aws_security_group.this.id
  from_port         = 8834
  to_port           = 8834
  ip_protocol       = "tcp"
  cidr_ipv4         = var.allowed_ip
}
