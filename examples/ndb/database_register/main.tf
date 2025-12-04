terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.8.0-beta.2"
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

## register PostgreSQL database with registered DBServer VM
resource "nutanix_ndb_register_database" "name" {
  database_type = "postgres_database"
  database_name=  "test-inst"
  description = "added by terraform"
  category =  "DEFAULT"

  // registered vm IP
  vm_ip = "{{ vm_ip }}"

  // optional 
  working_directory= "/tmp"
  
  reset_description_in_nx_cluster= false

  // time Machine Info
  time_machine_info {
    name= "test-pg-inst-regis"
    description= "description of tms"
    slaid=" {{ SLA ID}}"
    schedule {
      snapshottimeofday{
        hours= 16
        minutes= 0
        seconds= 0
      }
      continuousschedule{
        enabled=true
        logbackupinterval= 30
        snapshotsperday=1
      }
      weeklyschedule{
        enabled=true
        dayofweek= "WEDNESDAY"
      }
      monthlyschedule{
        enabled = true
        dayofmonth= "27"
      }
      quartelyschedule{
        enabled=true
        startmonth="JANUARY"
        dayofmonth= 27
      }
      yearlyschedule{
        enabled= false
        dayofmonth= 31
        month="DECEMBER"
      }
    }
  }
  postgress_info{

    // required args
    listener_port= "5432"
    db_password ="pass"
    db_name= "testdb1"

    // Optional with default values
    db_user= "postgres"
    backup_policy= "prefer_secondary"
    postgres_software_home= "{{ directory where the PostgreSQL database software is installed.}}"
    software_home= "{{ directory where the PostgreSQL database software is installed. }}"
   
  }
}


## register PostgreSQL database with instance not registered on VM
resource "nutanix_ndb_register_database" "name" {
  database_type = "postgres_database"
  database_name=  "test-inst"
  description = "added by terraform"
  category =  "DEFAULT"
  nx_cluster_id = "{{ cluster_ID }}"

  // registered vm info
  vm_ip = "{{ vm_ip }}"
  vm_username = "{{ vm_username }}"
  vm_password = "{{ vm_password }}"

  // optional 
  working_directory= "/tmp"
  
  reset_description_in_nx_cluster= false

  // time Machine Info
  time_machine_info {
    name= "test-pg-inst-regis"
    description= "description of tms"
    slaid=" {{ SLA ID}}"
    schedule {
      snapshottimeofday{
        hours= 16
        minutes= 0
        seconds= 0
      }
      continuousschedule{
        enabled=true
        logbackupinterval= 30
        snapshotsperday=1
      }
      weeklyschedule{
        enabled=true
        dayofweek= "WEDNESDAY"
      }
      monthlyschedule{
        enabled = true
        dayofmonth= "27"
      }
      quartelyschedule{
        enabled=true
        startmonth="JANUARY"
        dayofmonth= 27
      }
      yearlyschedule{
        enabled= false
        dayofmonth= 31
        month="DECEMBER"
      }
    }
  }
  postgress_info{

    // required args
    listener_port= "5432"
    db_password ="pass"
    db_name= "testdb1"

    // Optional with default values
    db_user= "postgres"
    backup_policy= "prefer_secondary"
    postgres_software_home= "{{ directory where the PostgreSQL database software is installed }}"
  }
}
