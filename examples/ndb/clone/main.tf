terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.8.0"
        }
    }
}

#defining nutanix configuration
provider "nutanix"{
  ndb_username = var.ndb_username
  ndb_password = var.ndb_password
  ndb_endpoint = var.ndb_endpoint
  insecure = true
}


## resource for cloning  using Point in time given time machine name

resource "nutanix_ndb_clone" "name" {
    time_machine_name = "test-pg-inst"
    name = "test-inst-tf-check"
    nx_cluster_id = "{{ nx_Cluster_id }}"
    ssh_public_key = "{{ sshkey }}"
    user_pitr_timestamp=  "{{ point_in_time }}"
    time_zone = "Asia/Calcutta"
    create_dbserver = true
    compute_profile_id = "{{ compute_profile_id }}"
    network_profile_id ="{{ network_profile_id }}"
    database_parameter_profile_id =  "{{ databse_profile_id }}"
    nodes{
        vm_name= "test_vm_clone"
        compute_profile_id = "{{ compute_profile_id }}"
        network_profile_id ="{{ network_profile_id }}"
        nx_cluster_id = "{{ nx_Cluster_id }}"
    }
    postgresql_info{
        vm_name="test_vm_clone"
        db_password= "pass"
    }
}

## resource for cloning  using snapshot given time machine name

resource "nutanix_ndb_clone" "name" {
    time_machine_name = "test-pg-inst"
    name = "test-inst-tf-check"
    nx_cluster_id = "{{ nx_Cluster_id }}"
    ssh_public_key = "{{ sshkey }}"
    snapshot_id=  "{{ snapshot_id }}"
    time_zone = "Asia/Calcutta"
    create_dbserver = true
    compute_profile_id = "{{ compute_profile_id }}"
    network_profile_id ="{{ network_profile_id }}"
    database_parameter_profile_id =  "{{ databse_profile_id }}"
    nodes{
        vm_name= "test_vm_clone"
        compute_profile_id = "{{ compute_profile_id }}"
        network_profile_id ="{{ network_profile_id }}"
        nx_cluster_id = "{{ nx_Cluster_id }}"
    }
    postgresql_info{
        vm_name="test_vm_clone"
        db_password= "pass"
    }
}

## resource for cloning with LCM Config with Registered DB Server VM

resource "nutanix_ndb_clone" "name" {
  time_machine_name = "test-pg-inst"
  name = "test-inst-tf-check"
  nx_cluster_id = "{{ nx_cluster_id }}"
  ssh_public_key = "{{ public_key }}"
  user_pitr_timestamp=  "{{ point in time }}"
  time_zone = "Asia/Calcutta"
  create_dbserver = false
  database_parameter_profile_id =  "{{ database_paramter_id }}"
  dbserver_id="{{ dbserver_id }}"
  nodes{
    vm_name="test-era-vm-regis"
    nx_cluster_id = "{{ nx_cluster_id }}"
    dbserver_id="{{ dbserver_id }}"
  }
  lcm_config{
    database_lcm_config{
      expiry_details{
        expire_in_days = 3
        expiry_date_timezone = "Asia/Calcutta"
      }
      refresh_details{
        refresh_in_days = 2
      }
    }
  }
  postgresql_info{
    vm_name="test-era-vm-regis"
    db_password= "pass"
  }
}