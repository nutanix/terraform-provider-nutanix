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


## resource to create maintenance window with Weekly as recurrence

resource nutanix_ndb_maintenance_window minWin{
    name = "test"
    description = "this is desc"
    recurrence = "WEEKLY"
    duration = 2
    day_of_week = "TUESDAY"
    start_time = "17:04:47" 
}


## resource to create maintenance window with Monthly as recurrence

resource nutanix_ndb_maintenance_window acctest-managed{
    name = "test"
    description = "this is desc"
    recurrence = "MONTHLY"
    duration = 2
    day_of_week = "TUESDAY"
    start_time = "17:04:47" 
    week_of_month= 4
}