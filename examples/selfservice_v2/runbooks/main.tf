terraform {
  required_providers {
    nutanix = {
      source  = "nutanixtemp/nutanix"
      version = "1.99.99"
    }
  }
}

provider "nutanix" {
  username = "admin"
  password = "Nutanix.123"
  endpoint = "10.44.76.58"
  insecure = true
  port     = 9440
}

resource "nutanix_calm_runbook" "TestRunbookCreate" {
  runbook_name = "TestRunbook"
  runbook_description = "Runbook description"
  project_uuid = "65b8817c-341f-4121-8faa-5e7bc94faa43"
  default_endpoint_name = "9987008909"
  task_list {
    task_name = "task1"
    task_type = "exec"
    task_script_type = "escript"
    task_script = "print(\"test script\")"
  }

  task_list {
    task_name = "task2"
    task_type = "exec"
    task_script_type = "escript"
    task_script = "print(\"test script\")"
  }
}

output "rb_state" {
  value = nutanix_calm_runbook.TestRunbookCreate.state
}

resource "nutanix_calm_runbook_execute" "TestRunbookExecute" {
  rb_name = nutanix_calm_runbook.TestRunbookCreate.runbook_name
}

output "rb_execution_state" {
  value = nutanix_calm_runbook_execute.TestRunbookExecute.state
}

resource "nutanix_calm_runbook_execute" "TestRunbookExecuteWithInputVariables" {
  rb_name = "TestRunbook"

  variable_list {
      name = "var1"
      value = "newval1"
    }

  variable_list {
      name = "var2"
      value = "10"
  }
}

resource "nutanix_calm_runbook_execute" "TestRunbookExecuteWithRunbookUUID" {
  rb_uuid = "ea66c5be-6bc1-dbf3-75d3-4c0c6568bdfb"
}
