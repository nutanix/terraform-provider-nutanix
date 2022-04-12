// resources/datasources used in this file were introduced in nutanix/nutanix version 1.5.0-beta.2
terraform {
    required_providers {
      nutanix = {
          source = "nutanix/nutanix"
          version = ">1.5.0-beta.2"
      }
    }
}

provider "nutanix" {
    username  = "user"
    password  = "pass"
    endpoint  = "10.x.xx.xx"
    insecure  = true
    port      = 9440
}


// resource to Create a new api key

resource "nutanix_foundation_central_api_keys" "apk"{
    alias = "test-fc"
}

// datasource to get an api key given its UUID
data "nutanix_foundation_central_api_keys" "k1"{
    key_uuid = resource.nutanix_foundation_central_api_keys.apk.key_uuid
}

output "k2"{
    value = data.nutanix_foundation_central_api_keys.k1
}

//datasource to List all the api keys created in Foundation Central.

data "nutanix_foundation_central_list_api_keys" "l1"{}

output "l2"{
    value = data.nutanix_foundation_central_list_api_keys.l1
}