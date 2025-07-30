# ------------------------------------------------------------------------------
# EC2 INSTANCE
# ------------------------------------------------------------------------------

resource "aws_instance" "this" {
  ami                    = var.ami_id
  user_data              = var.user_data
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.this.id]
  key_name               = aws_key_pair.this.key_name
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
  filename        = pathexpand("${var.key_directory}/${var.deployment_id}.pem")
  content         = tls_private_key.this.private_key_openssh
  file_permission = "400"
}

# ------------------------------------------------------------------------------
# SECURITY GROUP
# ------------------------------------------------------------------------------

resource "aws_security_group" "this" {}

resource "aws_vpc_security_group_egress_rule" "egress_any" {
  security_group_id = aws_security_group.this.id
  description       = "Allow all outbound traffic"

  ip_protocol = "-1"
  cidr_ipv4   = "0.0.0.0/0"
}

resource "aws_vpc_security_group_ingress_rule" "ingress_rules" {
  count             = length(var.ingress_rules)
  security_group_id = aws_security_group.this.id
  description       = var.ingress_rules[count.index].description

  ip_protocol    = "tcp"
  prefix_list_id = var.ingress_rules[count.index].prefix_list_id
  cidr_ipv4      = var.ingress_rules[count.index].cidr_ipv4
  from_port      = var.ingress_rules[count.index].port
  to_port        = var.ingress_rules[count.index].port
}
