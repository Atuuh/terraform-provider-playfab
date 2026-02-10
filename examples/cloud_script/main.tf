terraform {
  required_providers {
    playfab = {
      source = "registry.terraform.io/atuuh/playfab"
    }
  }
}

provider "playfab" {
  title_id   = "127DE0"
  secret_key = "S66ZQ7FKD7N3BCF4SPDJIFZNYTKMUWNZOJRYNWXO5MYZ6KFR3O"
}

data "playfab_cloud_script" "example" {}

output "example_functions" {
  value = data.playfab_cloud_script.example
}

resource "playfab_function" "test" {
  name = "Frank's Fun-ction"
  url = "http://0.0.0.0:12345/api/dance/party"
  trigger_type = "http"
}

output "test_function" {
  value = playfab_function.test
}