resource "docker_image" "foo" {
  name         = "nginx:latest"
  keep_locally = true
}

resource "docker_container" "foo" {
  name  = "tf-test"
  image = docker_image.foo.latest

  upload {
    source     = "%s"
    file       = "/terraform/test.txt"
    executable = true
  }
}
