nutanix_central_endpoint = "nconprem-10-44-76-35.ccpnx.com"
tenant_password         = "nutanix/4u"
tenant_username         = "apple_user_101@apple.lab.local"
tenant_name              = "cpaas-iac-compute-tenant"
domain_name              = "goku"

ad_user_query        = "apple_user_011"
ad_group_query       = "apple_user_group_7"

# Directory Service Configuration
directory_service_name = "tenant-ad"
directory_service_url = "ldap://10.44.76.35:389"
directory_service_username = "apple_user_101@apple.lab.local"
directory_service_password = "nutanix/4u"
directory_service_domain_name = "apple.lab.local"
whitelisted_groups = ["apple_user_group_7"]