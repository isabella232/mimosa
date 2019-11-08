# ----------------------------------------------------------------
#
# Deploy AWS instances
#
# Code assumes that a key pair called "mimosa" exists in the region already
#
# Instance lifetime is 1 day but can be changed below
#
# ----------------------------------------------------------------

# ----------------------------------------------------------------
#
# DEPLOY - Specify the region/total of your choice:
#
# terraform apply -auto-approve -var total=2 -var awsregion=eu-west-1 -state=eu-west-1.tfstate
#
# MAKE SURE YOU SET THE STATE FILE TO MATCH YOUR REGION!
#
# ----------------------------------------------------------------

# ----------------------------------------------------------------
#
# TEAR DOWN - Specify the region/total of your choice:
#
# terraform destroy -auto-approve -var total=2 -var awsregion=eu-west-1 -state=eu-west-1.tfstate
#
# TRULY, MAKE SURE YOU SET THE STATE FILE TO MATCH YOUR REGION!
#
# ----------------------------------------------------------------

variable "awsregion" {
  type = string
}

variable "total" {
  type = number
}

provider "aws" {
  region = var.awsregion
}

resource "aws_instance" "ubuntu1804" {

  count = var.total

  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.nano"
  key_name      = "mimosa"

  tags = {
    lifetime = "1d"
    Name     = "mimosa-ubuntu1804-${count.index}"
  }

}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}
