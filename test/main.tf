terraform {
    backend "local" {
    path = "./terraform.tfstate"
  }
}

resource "local_file" "test" {
  content  = "foo!"
  filename = "${path.module}/test.txt"
}
