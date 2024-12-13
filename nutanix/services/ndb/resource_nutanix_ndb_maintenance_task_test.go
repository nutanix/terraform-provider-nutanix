package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceMaintenaceTaskName = "nutanix_ndb_maintenance_task.acctest-managed"

func TestAccEra_MaintenanceTask(t *testing.T) {
	name := "test-maintenance-acc"
	desc := "this is desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraMaintenanceTask(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "name", name),
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "description", desc),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.#"),
					resource.TestCheckResourceAttr(resourceMaintenaceTaskName, "entity_task_association.0.entity_type", "ERA_DBSERVER"),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.1.task_type"),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.0.task_type"),
				),
			},
		},
	})
}

func TestAccEra_MaintenanceTask_Update(t *testing.T) {
	name := "test-maintenance-acc"
	desc := "this is desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraMaintenanceTask(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "name", name),
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "description", desc),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.#"),
					resource.TestCheckResourceAttr(resourceMaintenaceTaskName, "entity_task_association.0.entity_type", "ERA_DBSERVER"),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.1.task_type"),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.0.task_type"),
				),
			},
			{
				Config: testAccEraMaintenanceTaskUpdate(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "name", name),
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "description", desc),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.#"),
					resource.TestCheckResourceAttr(resourceMaintenaceTaskName, "entity_task_association.#", "1"),
					resource.TestCheckResourceAttr(resourceMaintenaceTaskName, "entity_task_association.0.entity_type", "ERA_DBSERVER"),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.0.task_type"),
				),
			},
		},
	})
}

func TestAccEra_MaintenanceTask_UpdateWithNoTask(t *testing.T) {
	name := "test-maintenance-acc"
	desc := "this is desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraMaintenanceTask(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "name", name),
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "description", desc),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.#"),
					resource.TestCheckResourceAttr(resourceMaintenaceTaskName, "entity_task_association.0.entity_type", "ERA_DBSERVER"),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.1.task_type"),
					resource.TestCheckResourceAttrSet(resourceMaintenaceTaskName, "entity_task_association.0.task_type"),
				),
			},
			{
				Config: testAccEraMaintenanceTaskUpdateWithNoTask(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "name", name),
					resource.TestCheckResourceAttr(resourceMaintenaceWindowName, "description", desc),
					resource.TestCheckResourceAttr(resourceMaintenaceTaskName, "entity_task_association.#", "0"),
				),
			},
		},
	})
}

func testAccEraMaintenanceTask(name, desc string) string {
	return fmt.Sprintf(`
		resource nutanix_ndb_maintenance_window acctest-managed{
			name = "%[1]s"
			description = "%[2]s"
			recurrence = "WEEKLY"
			duration = 2
			day_of_week = "TUESDAY"
			start_time = "17:04:47" 
		}

		data "nutanix_ndb_dbservers" "dbservers"{}

		resource nutanix_ndb_maintenance_task acctest-managed{
			dbserver_id = [
				data.nutanix_ndb_dbservers.dbservers.dbservers.0.id
			]
			maintenance_window_id = resource.nutanix_ndb_maintenance_window.acctest-managed.id
			tasks{
			  task_type = "OS_PATCHING"
			}
			tasks {
			  task_type = "DB_PATCHING"
			}
		}
	`, name, desc)
}

func testAccEraMaintenanceTaskUpdate(name, desc string) string {
	return fmt.Sprintf(`
		resource nutanix_ndb_maintenance_window acctest-managed{
			name = "%[1]s"
			description = "%[2]s"
			recurrence = "WEEKLY"
			duration = 2
			day_of_week = "TUESDAY"
			start_time = "17:04:47" 
		}

		data "nutanix_ndb_dbservers" "dbservers"{}

		resource nutanix_ndb_maintenance_task acctest-managed{
			dbserver_id = [
				data.nutanix_ndb_dbservers.dbservers.dbservers.0.id
			]
			maintenance_window_id = resource.nutanix_ndb_maintenance_window.acctest-managed.id
			tasks {
			  task_type = "DB_PATCHING"
			}
		}
	`, name, desc)
}

func testAccEraMaintenanceTaskUpdateWithNoTask(name, desc string) string {
	return fmt.Sprintf(`
		resource nutanix_ndb_maintenance_window acctest-managed{
			name = "%[1]s"
			description = "%[2]s"
			recurrence = "WEEKLY"
			duration = 2
			day_of_week = "TUESDAY"
			start_time = "17:04:47" 
		}

		data "nutanix_ndb_dbservers" "dbservers"{}

		resource nutanix_ndb_maintenance_task acctest-managed{
			dbserver_id = [
				data.nutanix_ndb_dbservers.dbservers.dbservers.0.id
			]
			maintenance_window_id = resource.nutanix_ndb_maintenance_window.acctest-managed.id
		}
	`, name, desc)
}
