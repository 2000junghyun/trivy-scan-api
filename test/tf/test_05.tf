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