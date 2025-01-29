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
  runbook_name = "Jaldi_run_hoja_ghar_jaana_hai"
  runbook_description = "Runbook description"
  project_uuid = "65b8817c-341f-4121-8faa-5e7bc94faa43"
  default_endpoint_name = "9987008909"
  task_list {
    task_name = "task1"
    task_type = "exec"
    task_script_type = "escript"
    task_script = "print(\"hi\")"
  }

  task_list {
    task_name = "task2"
    task_type = "exec"
    task_script_type = "escript"
    task_script = "print(\"Ba bye\")"
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
