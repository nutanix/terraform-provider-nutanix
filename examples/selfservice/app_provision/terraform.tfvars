#define values to the variables to be used in terraform file
nutanix_username = "admin"
nutanix_password = "password"
nutanix_endpoint = "10.xx.xx.xx"
nutanix_port = 9440

blueprint_name = "name_of_blueprint"
app_name = "name_of_app"
app_description = "description_of_app"
file_name = "runtime_value.json" // any valid json filename to dump runtime editables fetched
substrate_name = "name_of_substrate" // name of substrate whose value you want to change at runtime
system_action_name = "stop" // valid actions are ["start", "stop", "restart"]
