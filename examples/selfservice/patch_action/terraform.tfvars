#define values to the variables to be used in terraform file
nutanix_username = "admin"
nutanix_password = "password"
nutanix_endpoint = "10.xx.xx.xx"
nutanix_port = 9440

blueprint_name = "name_of_blueprint"
app_name = "name_of_app"
app_description = "description_of_app"
patch_name = "name_of_patch_action"
config_name = "name_of_config" // same as patch_name
memory_size_mib = 1024 // integer value (size in Mib)
num_sockets = 2
num_vcpus_per_socket = 2
category_value = "Key:Value" // (e.g "AppType:Default")
add_operation = "add"
delete_operation = "delete"
disk_size_mib = 3072 // integer value (size in Mib)
index = "0" // dummy value of index
subnet_uuid = "1234-5678-9012" // valid subnet uuid present in project assigned to application