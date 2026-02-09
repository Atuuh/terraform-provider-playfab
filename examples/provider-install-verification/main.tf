terraform {
  required_providers {
    playfab = {
      source = "registry.terraform.io/atuuh/playfab"
    }
  }
}

provider "playfab" {}

data "playfab_coffees" "example" {}
