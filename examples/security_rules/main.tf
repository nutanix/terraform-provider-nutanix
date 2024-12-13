###################################
#            General
###################################

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.3.0"
    }
  }
}

#define nutanix provider configuration
provider "nutanix" {
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_endpoint
  port         = var.nutanix_port
  insecure     = true
  wait_timeout = 60
}

###################################
#            Data input
###################################

#retrieve the cluster corresponding to the cluster name variable
data "nutanix_cluster" "cluster" {
  name = var.cluster_name
}
#retrieve the network corresponding to the network name variable
data "nutanix_subnet" "subnet" {
  subnet_name = var.subnet_name
}

#retrieve image details
data "nutanix_image" "image"{
  image_name = "CentOS-7 Generic Cloud managed by Terraform"
}

###################################
#            Categories
###################################
#create category named DEMO-DOLIBARR-TIER
resource "nutanix_category_key" "demo-dolibarr-tier" {
  name        = "DEMO-DOLIBARR-TIER"
  description = "TIER Category Key"
}

#add category values in above created category
resource "nutanix_category_value" "LB" {
  name        = nutanix_category_key.demo-dolibarr-tier.id
  description = "LoadBalancer Tier"
  value       = "LoadBalancer"
}
resource "nutanix_category_value" "WEB" {
  name        = nutanix_category_key.demo-dolibarr-tier.id
  description = "Webserver Tier"
  value       = "Webserver"
}
resource "nutanix_category_value" "DB" {
  name        = nutanix_category_key.demo-dolibarr-tier.id
  description = "Database Tier"
  value       = "Database"
}

###################################
#            uSeg
###################################
#creating security rules

resource "nutanix_network_security_rule" "loadbalancer" {
  name            = "DEMO-TERRAFORM-DOLIBARR"               #name of the uSeg rules
  description     = "DEMO-TERRAFORM-DOLIBARR"
  app_rule_action = "MONITOR"
 

  app_rule_target_group_peer_specification_type = "FILTER"
  app_rule_target_group_default_internal_policy = "ALLOW_ALL"
  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"
  app_rule_target_group_filter_kind_list = [
    "vm"
  ]

  app_rule_target_group_filter_params {
    name = "AppType"
    values = [
      "Default"
    ]
  }

  #filter with previously created category key-value
  app_rule_target_group_filter_params {
    name = "DEMO-DOLIBARR-TIER"
    values = [
      "LoadBalancer"
    ]
  }

  app_rule_inbound_allow_list {
    ip_subnet               = "10.xx.xx.xx"
    ip_subnet_prefix_length = "0"
    peer_specification_type = "IP_SUBNET"
    protocol                = "TCP"
    tcp_port_range_list {
      end_port   = 88
      start_port = 88
    }
    tcp_port_range_list {
      end_port   = 89
      start_port = 89
    }
  }

  app_rule_outbound_allow_list {
    ip_subnet               = "10.xx.xx.xx"
    ip_subnet_prefix_length = "0"
    peer_specification_type = "IP_SUBNET"
    protocol                = "TCP"
    tcp_port_range_list {
      end_port   = 88
      start_port = 88
    }
    tcp_port_range_list {
      end_port   = 89
      start_port = 89
    }
  }
}

resource "nutanix_network_security_rule" "webserver" {
  name            = "DEMO-TERRAFORM-DOLIBARR"
  description     = "DEMO-TERRAFORM-DOLIBARR"
  app_rule_action = "MONITOR"
 

  app_rule_target_group_peer_specification_type = "FILTER"
  app_rule_target_group_default_internal_policy = "ALLOW_ALL"
  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"
  app_rule_target_group_filter_kind_list = [
    "vm"
  ]

  app_rule_target_group_filter_params {
    name = "AppType"
    values = [
      "Default"
    ]
  }
  #filter with previously created category key-value
  app_rule_target_group_filter_params {
    name = "DEMO-DOLIBARR-TIER"
    values = [
      "Webserver"
    ]
  }
}

resource "nutanix_network_security_rule" "database" {
  name            = "DEMO-TERRAFORM-DOLIBARR"
  description     = "DEMO-TERRAFORM-DOLIBARR"
  app_rule_action = "MONITOR"
 

  app_rule_target_group_peer_specification_type = "FILTER"
  app_rule_target_group_default_internal_policy = "ALLOW_ALL"
  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"
  app_rule_target_group_filter_kind_list = [
    "vm"
  ]

  app_rule_target_group_filter_params {
    name = "AppType"
    values = [
      "Default"
    ]
  }
  #filter with previously created category key-value
  app_rule_target_group_filter_params {
    name = "DEMO-DOLIBARR-TIER"
    values = [
       "Database"
    ]
  }
}

resource "nutanix_network_security_rule" "loadbalancer-webserver" {
  name            = "DEMO-TERRAFORM-DOLIBARR"
  description     = "DEMO-TERRAFORM-DOLIBARR"
  app_rule_action = "MONITOR"
 

  app_rule_inbound_allow_list {
    filter_type = "CATEGORIES_MATCH_ALL"
    filter_params {
      name = "AppType"
      values = [
        "Default"
      ]
    }    
    #filter with previously created category key-value
    filter_params {
      name = "DEMO-DOLIBARR-TIER"
      values = [
        "LoadBalancer"
      ]
    }
    filter_kind_list        = ["vm"]
    peer_specification_type = "FILTER"
    protocol                = "TCP"
    tcp_port_range_list {
      end_port = 88
      start_port = 88
    }
  }


  app_rule_target_group_peer_specification_type = "FILTER"
  app_rule_target_group_default_internal_policy = "ALLOW_ALL"
  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"
  app_rule_target_group_filter_kind_list = [
    "vm"
  ]

  app_rule_target_group_filter_params {
    name = "AppType"
    values = [
      "Default"
    ]
  }
  #filter with previously created category key-value
  app_rule_target_group_filter_params {
    name = "DEMO-DOLIBARR-TIER"
    values = [
      "Webserver"
    ]
  }
}

#create security rule
resource "nutanix_network_security_rule" "webserver-database" {
  name            = "DEMO-TERRAFORM-DOLIBARR"
  description     = "DEMO-TERRAFORM-DOLIBARR"
  app_rule_action = "MONITOR"
 

  app_rule_inbound_allow_list {
    filter_type = "CATEGORIES_MATCH_ALL"
    filter_params {
      name = "AppType"
      values = [
        "Default"
      ]
    }    

    #filter with previously created category key-value
    filter_params {
      name = "DEMO-DOLIBARR-TIER"
      values = [
        "Webserver"
      ]
    }
    filter_kind_list        = ["vm"]
    peer_specification_type = "FILTER"
    protocol                = "TCP"
    tcp_port_range_list {
      end_port = 3333
      start_port = 3333
    }
  }


  app_rule_target_group_peer_specification_type = "FILTER"
  app_rule_target_group_default_internal_policy = "ALLOW_ALL"
  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"
  app_rule_target_group_filter_kind_list = [
    "vm"
  ]

  app_rule_target_group_filter_params {
    name = "AppType"
    values = [
      "Default"
    ]
  }
  app_rule_target_group_filter_params {
    name = "DEMO-DOLIBARR-TIER"
    values = [
      "Database"
    ]
  }
}
