module "rules_two" {
  source = "terraform-aws-modules/security-group/aws"

  create_sg         = false
  vpc_id            = "SECRET"
  security_group_id = "SECRET"

  ingress_with_cidr_blocks = [
    {
      from_port   = 22
      to_port     = 22
      protocol    = "-1"
      cidr_blocks = "0.0.0.0/0"
    },
  ]
}

resource "aws_security_group" "test2" {
  name        = "allow_ssh_from_world"
  vpc_id      = "SECRET"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}