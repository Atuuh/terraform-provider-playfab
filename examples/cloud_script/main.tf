terraform {
  required_providers {
    playfab = {
      source = "registry.terraform.io/atuuh/playfab"
    }
  }
}

variable "playfab_title_id" {
  type = string
}

variable "playfab_secret_key" {
  type      = string
  sensitive = true
}

provider "playfab" {
  title_id   = var.playfab_title_id
  secret_key = var.playfab_secret_key
}

data "playfab_cloud_script" "example" {}

output "example_functions" {
  value = data.playfab_cloud_script.example
}

resource "playfab_function" "test" {
  name         = "Frank's Fun-ction"
  url          = "http://0.0.0.0:12345/api/dance/party"
  trigger_type = "http"
}

output "test_function" {
  value = playfab_function.test
}
