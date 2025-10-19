resource "aws_launch_template" "web" {
  name_prefix   = "${var.project}-web-lt-"
  image_id      = var.ami_id
  instance_type = var.instance_type
  key_name      = var.key_pair_name

  user_data = base64encode(file(var.user_data_path))

  network_interfaces {
    associate_public_ip_address = false
    security_groups             = [var.web_sg_id]
  }

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "${var.project}-web"
    }
  }
}

resource "aws_autoscaling_group" "web" {
  name                      = "${var.project}-web-asg"
  max_size                  = 3
  min_size                  = 1
  desired_capacity          = 1
  vpc_zone_identifier       = var.web_subnet_ids
  target_group_arns         = [var.web_target_group_arn]

  launch_template {
    id      = aws_launch_template.web.id
    version = "$Latest"
  }

  tag {
    key                 = "Name"
    value               = "${var.project}-web"
    propagate_at_launch = true
  }
}