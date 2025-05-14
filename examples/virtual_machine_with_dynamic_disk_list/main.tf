#Here we will create a vm using an existing image and add dyanmic disk blocks as per 
#the variable "disk_sizes" present in terraform.tfvars file.
#Note - Replace appropriate values of variables in terraform.tfvars file as per setup

terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.3.0"
        }
    }
}

#defining nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = var.nutanix_port
  insecure = true
}

#pull existing image data (can upload image as well using nutanix_image resource)
data "nutanix_image" "centos"{
  image_name = "Centos7-Base"
}

#pull desired cluster data
data "nutanix_cluster" "cluster"{
  name = var.cluster_name
}

#pull desired subnet data
data "nutanix_subnet" "subnet"{
  subnet_name = var.subnet_name
}

#create a virtual machine
resource "nutanix_virtual_machine" "dev-vm-demo" {
  count = 1
  name = "test-vm-2937479"
  cluster_uuid = data.nutanix_cluster.cluster.id
  num_vcpus_per_socket = "1"
  num_sockets = "4"
  memory_size_mib = 4096

  #add basic disk with centos image
  disk_list {
    data_source_reference = {
      kind = "image"
      uuid = data.nutanix_image.centos.id
    }
  }

  #add nic
  nic_list {
    subnet_uuid = data.nutanix_subnet.subnet.id
  }

  #dynamically add disk blocks of various sizes as per variable list - disk_sizes
  dynamic "disk_list" {
    for_each = [for disk in var.disk_sizes : disk]
    content {
      disk_size_mib = disk_list.value
    }
  }
}