#replace the values as per setup configuration
nutanix_username = "admin"
nutanix_password = "Nutanix/123456"
nutanix_endpoint = "10.xx.xx.xx"
nutanix_port     = 9440

#replace this values as per the setup
vm_uuid = "<vm-uuid>"

#this variable will be used in adding disks to vm in main.tf
disk_sizes = [1024, 1024, 2048]
