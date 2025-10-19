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