terraform {
  required_providers {
    nutanix = {
      source  = "nutanixtemp/nutanix"
      version = "1.99.99"
    }
  }
}

provider "nutanix" {
  username = "admin"
  password = "Nutanix.123"
  endpoint = "10.44.76.58"
  insecure = true
  port     = 9440
}
resource "nutanix_calm_runbook_execute" "TestRunbook" {
  rb_name = "rb289989"

  variable_list {
      name = "var1"
      value = "newval1"
    }

  variable_list {
      name = "var2"
      value = "10"
  }
}

resource "nutanix_calm_runbook_execute" "TestRunbook2" {
  rb_name = "rbsimple"
}

resource "nutanix_calm_runbook_execute" "TestRunbook3" {
  rb_uuid = "ea66c5be-6bc1-dbf3-75d3-4c0c6568bdfb"
}

