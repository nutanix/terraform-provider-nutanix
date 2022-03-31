// resources/datasources used in this file were introduced in nutanix/nutanix version >=1.4.2
terraform{
    required_providers{
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.4.2"
        }
    }
}

// default foundation_port is 8000 so can be ignored
provider "nutanix" {
    // foundation_port = 8000
    foundation_endpoint = "10.xx.xx.xx"
}

/*
Description:
* Here we will be uploading
 - 1 aos image with name nos-tempfile5.20.1.1.tar
 - 1 esx hypervisor image with name esx_image.iso
 - 1 hyperv hypervisor image with name xen_image
* Source -> this field will always be the fie path in local setup where this terraform file is run.
Note : Only .tar for aos and .iso for hypervisor are displayed by data sources of them, so avaid uploading invalid file types.
*/


// upload aos image
resource "nutanix_foundation_image" "image1" {
  source = "nutanix_installer-x86_64.tar"
  filename = "nos-tempfile5.20.1.1.tar"
  installer_type = "nos"
}

// upload esx hypervisor image
resource "nutanix_foundation_image" "image2" {
  source = "../../../files/VMware-VMvisor-Installer.x86_64.iso"
  filename = "esx_image.iso"
  installer_type = "esx"
}

// upload xen hypervisor image
resource "nutanix_foundation_image" "image3" {
  source = "../../../files/Xen-installer.x86_64.iso"
  filename = "xen_image.iso"
  installer_type = "xen"
}

// Fetch all aos packages detais once upload finishes
data "nutanix_foundation_nos_packages" "nos" {
  depends_on = [resource.nutanix_foundation_image.image1]
}

// Fetch all aos hypervisor image details once upload finishes
data "nutanix_foundation_hypervisor_isos" "hyper" {
  depends_on  = [resource.nutanix_foundation_image.image2, resource.nutanix_foundation_image.image3]
}

output "nos" {
  value = data.nutanix_foundation_nos_packages.nos
}
output "hypervisors" {
  value = data.nutanix_foundation_hypervisor_isos.hyper
}