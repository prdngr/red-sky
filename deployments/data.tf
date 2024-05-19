data "aws_ami" "nessus" {
  most_recent = true

  filter {
    name   = "name"
    values = ["Nessus *"]
  }

  filter {
    name   = "product-code"
    values = ["8fn69npzmbzcs4blc4583jd0y"]
  }

  owners = ["679593333241"]
}
