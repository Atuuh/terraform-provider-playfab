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
