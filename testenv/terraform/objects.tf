
// NOTE: The object-store VLAN-800 subnet (objects.800) is NOT created here.
// prepare_pc.py owns it: ensure_flow_managed_subnet() creates-or-reuses the shared
// VLAN-800 managed subnet from object_store.subnet.* (see config.yaml), and the Flow
// Controller reuses that same subnet (AHV basic networking allows only one subnet per
// VLAN id). Defining it in Terraform too would create two identical subnets / collide,
// so the resource lives solely in prepare_pc.py. This keeps a single source of truth.

# locals {
#   object_subnet = nutanix_subnet_v2.object-subnet.name
# }



