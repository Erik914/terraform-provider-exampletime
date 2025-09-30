terraform {
  required_providers {
    exampletime = {
      source = "hashicorp.com/edu/exampletime"
    }
  }
}

provider "exampletime" {}

output "label" {
  value = provider::exampletime::label_split("test1/test2/test3")
}
